package service

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/preconfigs"
	"github.com/miaokobot/miaospeed/utils"
	"github.com/miaokobot/miaospeed/utils/structs"

	"github.com/miaokobot/miaospeed/service/matrices"
	"github.com/miaokobot/miaospeed/service/taskpoll"
)

type WsHandler struct {
	Serve func(http.ResponseWriter, *http.Request)
}

func (wh *WsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if wh.Serve != nil {
		wh.Serve(rw, r)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func InitServer() {
	if utils.GCFG.Binder == "" {
		utils.DErrorf("MiaoSpeed Server | Cannot listening the binder, bind=%s", utils.GCFG.Binder)
		os.Exit(1)
	}

	utils.DInfof("MiaoSpeed Server | Start Listening, bind=%s", utils.GCFG.Binder)

	wsHandler := WsHandler{
		Serve: func(rw http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(rw, r, nil)
			if err != nil {
				utils.DErrorf("MiaoServer Test | Socket establishing error, error=%s", err.Error())
				return
			}
			defer conn.Close()

			var poll *taskpoll.TaskPollController

			batches := structs.NewAsyncMap[string, bool]()
			cancel := func() {
				if poll != nil {
					for id := range batches.ForEach() {
						poll.Remove(id, taskpoll.TPExitInterrupt)
					}
				}
			}

			defer cancel()
			for {
				sr := interfaces.SlaveRequest{}
				err := conn.ReadJSON(&sr)
				if err != nil {
					if !strings.Contains(err.Error(), "EOF") && !strings.Contains(err.Error(), "reset by peer") {
						utils.DErrorf("MiaoServer Test | Task receiving error, error=%s", err.Error())
					}

					return
				}

				verified := utils.GCFG.VerifyRequest(&sr)
				utils.DLogf("MiaoServer Test | Receive Task, name=%s invoker=%v matrices=%v payload=%d verify=%v", sr.Basics.ID, sr.Basics.Invoker, sr.Options.Matrices, len(sr.Nodes), verified)

				// verify token
				if !verified {
					conn.WriteJSON(&interfaces.SlaveResponse{
						Error: "cannot verify the request, please check your token",
					})
					return
				}
				sr.Challenge = ""

				// verify invoker
				if !utils.GCFG.InWhiteList(sr.Basics.Invoker) {
					conn.WriteJSON(&interfaces.SlaveResponse{
						Error: "the bot id is not in the whitelist",
					})
					return
				}

				// find all matrices
				matrices := matrices.FindBatchFromEntry(sr.Options.Matrices)

				// extra macro from the matrices
				macros := ExtractMacrosFromMatrices(matrices)

				// select poll
				if structs.Contains(macros, interfaces.MacroSpeed) {
					if utils.GCFG.NoSpeedFlag {
						conn.WriteJSON(&interfaces.SlaveResponse{
							Error: "speedtest is disabled on backend",
						})
						return
					}
					poll = SpeedTaskPoll
				} else {
					poll = ConnTaskPoll
				}
				utils.DLogf("MiaoServer Test | Receive Task, name=%s poll=%s", sr.Basics.ID, poll.Name())

				// build testing item
				item := poll.Push((&TestingPollItem{
					id:       utils.RandomUUID(),
					name:     sr.Basics.ID,
					request:  &sr,
					matrices: sr.Options.Matrices,
					macros:   macros,
					onProcess: func(self *TestingPollItem, idx int, result interfaces.SlaveEntrySlot) {
						conn.WriteJSON(&interfaces.SlaveResponse{
							ID:               self.ID(),
							MiaoSpeedVersion: utils.VERSION,
							Progress: &interfaces.SlaveProgress{
								Record:  result,
								Index:   idx,
								Queuing: poll.AwaitingCount(),
							},
						})
					},
					onExit: func(self *TestingPollItem, exitCode taskpoll.TaskPollExitCode) {
						batches.Del(self.ID())
						conn.WriteJSON(&interfaces.SlaveResponse{
							ID:               self.ID(),
							MiaoSpeedVersion: utils.VERSION,
							Result: &interfaces.SlaveTask{
								Request: sr,
								Results: self.results.ForEach(),
							},
						})
					},
				}).Init())

				batches.Set(item.ID(), true)
			}
		},
	}

	server := http.Server{
		Handler:   &wsHandler,
		TLSConfig: preconfigs.MakeSelfSignedTLSServer(),
	}

	if strings.HasPrefix(utils.GCFG.Binder, "/") {
		unixListener, err := net.Listen("unix", utils.GCFG.Binder)
		if err != nil {
			utils.DErrorf("MiaoServer Launch | Cannot listen on unixsocket %s, error=%s", utils.GCFG.Binder, err.Error())
			os.Exit(1)
		}
		server.Serve(unixListener)
	} else {
		netListener, err := net.Listen("tcp", utils.GCFG.Binder)
		if err != nil {
			utils.DErrorf("MiaoServer Launch | Cannot listen on socket %s, error=%s", utils.GCFG.Binder, err.Error())
			os.Exit(1)
		}
		if utils.GCFG.MiaoKoSignedTLS {
			server.ServeTLS(netListener, "", "")
		} else {
			server.Serve(netListener)
		}

	}
}

func CleanUpServer() {
	if strings.HasPrefix(utils.GCFG.Binder, "/") {
		os.Remove(utils.GCFG.Binder)
	}
}
