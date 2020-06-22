import React from 'react'
import './App.css'

import { JSONDecoder, RunningQuery, SoyuzClient } from 'soyuz'
import { schema } from './schema'
import { DialWebsocketClient } from './websocketClient'
import { CodeViewer } from './CodeViewer'

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
        <CodeViewer
            language="json"
            data={JSON.stringify(this.state, undefined, '\t')}
        />
      </div>
    )
  }

  private async startClient() {
    // TODO: connect websocket
    // TODO: re-construct soyuz client is websocket reconnects
    try {
      this.soyuzClient = await DialWebsocketClient(
        'ws://localhost:8093/ws',
        schema
      )
    } catch (e) {
      /* tslint:disable-next-line */
      console.error('dial websocket client', e)
      return
    }
    if (!this.soyuzClient) {
      return
    }
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
