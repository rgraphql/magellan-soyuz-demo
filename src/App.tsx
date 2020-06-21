import React from 'react'
import logo from './logo.svg'
import './App.css'

/*
client = new SoyuzClient(schema, (msg: rgraphql.IRGQLClientMessage) => {
    if (!sender) {
        sender = this.rpc.startCall(
            {
                rgqlClientMessage: msg,
                rpcId: ipc.RPC.RPC_RGQLClientMessage
            },
            (rmsg: ipc.IRPCMessage) => {
                const smsg = rmsg.rgqlServerMessage
                if (!smsg) {
                    return
                }
                client.handleMessages([smsg])
            }
        )
        return
    }

    sender({
        rgqlClientMessage: msg,
        rpcId: ipc.RPC.RPC_RGQLClientMessage
    })
})
*/

class App extends React.Component {
  public render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.tsx</code> and save to reload.
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
      </div>
    )
  }
}

export default App
