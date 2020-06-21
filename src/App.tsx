import React from 'react'
import { rgraphql } from 'rgraphql'
import './App.css'

import { JSONDecoder, RunningQuery, SoyuzClient } from 'soyuz'

import { schema } from './schema'

const AppDemoQuery = `{
  counter
  names
  allPeople {
    name
    height
  }
}`

interface AppProps {}

interface AppState {
  counter?: number
  names?: string[]
  allPeople?: any[]
}

class App extends React.Component<AppProps, AppState> {
  private soyuzClient?: SoyuzClient
  private query?: RunningQuery

  constructor(params: {}) {
    super(params)
    this.state = {}
  }

  public componentWillMount() {
      this.startClient()
  }

  public render() {
    return (
      <div className="App">
        <header className="App-header">
          <p>
            <code>{JSON.stringify(this.state, undefined, '\t')}</code>
          </p>
        </header>
      </div>
    )
  }

    private startClient() {
        // TODO: connect websocket
        // TODO: re-construct soyuz client is websocket reconnects
        this.soyuzClient = new SoyuzClient(
            schema,
            (msg: rgraphql.IRGQLClientMessage) => {
                // Transmit the message to the server.
                // TODO
                /* tslint:disable-next-line */
                console.log('Transmitting message to server:', msg)
            }
        )
        this.query = this.soyuzClient.parseQuery(AppDemoQuery)
        this.query.attachHandler(
            new JSONDecoder(
                this.soyuzClient.getQueryTree().getRoot(),
                this.query.getQuery(),
                (val: any) => {
                    if (val) {
                        this.setState(val)
                    }
                }
            )
        )
    }
}

export default App
