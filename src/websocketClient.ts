import * as graphql from 'graphql'
import { rgraphql } from 'rgraphql'
import { demo } from './pb/demo'
import { SoyuzClient } from 'soyuz'

// DialWebsocketClient dials a SoyuzClient over websocket.
export async function DialWebsocketClient(
  serverURL: string,
  schema: graphql.GraphQLSchema
): Promise<SoyuzClient> {
  const ws = new WebSocket(serverURL)
  return new Promise<SoyuzClient>((resolve, reject) => {
    let client: SoyuzClient
    ws.onopen = (e: Event) => {
      console.log('connected')
      client = new SoyuzClient(schema, (msg: rgraphql.IRGQLClientMessage) => {
        // Transmit the message to the server.
        /* tslint:disable-next-line */
        console.log('tx:', msg)
        const data = demo.RPCMessage.encode({
          rpcId: demo.RPC.RPC_RGQLClientMessage,
          rgqlClientMessage: msg,
        }).finish()
        ws.send(data)
      })
      resolve(client)
    }
    ws.onclose = (e: Event) => {
      console.log('disconnected')
      if (!client) {
        reject(new Error('connection failed'))
      }
    }
    ws.onerror = (err: Event) => {
      if (!client) {
        reject(err)
      }
    }
    ws.onmessage = async (e: MessageEvent) => {
      let data: Uint8Array
      const eventData = e.data
      if (eventData instanceof Uint8Array) {
        data = eventData
      } else {
        const dataBlob: Blob = e.data
        const dataArrayBuffer = await new Response(dataBlob).arrayBuffer()
        data = new Uint8Array(dataArrayBuffer)
      }
      try {
        const msg = demo.RPCMessage.decode(data)
        switch (msg.rpcId) {
          case demo.RPC.RPC_RGQLServerMessage:
            if (msg.rgqlServerMessage) {
              client.handleMessages([msg.rgqlServerMessage])
            }
            break
          case demo.RPC.RPC_Ping:
            break
          default:
            /* tslint:disable-next-line */
            console.error('unhandled rpc type', msg.rpcId)
        }
      } catch (e) {
        /* tslint:disable-next-line */
        console.error('handle message', e)
      }
    }
  })
}
