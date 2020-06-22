import * as monacoEditor from 'monaco-editor/esm/vs/editor/editor.api'
import * as React from 'react'
import MonacoEditor from 'react-monaco-editor'

import './CodeViewer.css'

// ICodeViewerProps are the object container props.
export interface ICodeViewerProps {
  // language is the data language
  language?: string
  // data is the initial data to show
  data?: string
}

// ICodeViewerState is the query overlay state.
export interface ICodeViewerState {}

const editorOptions: monacoEditor.editor.IEditorConstructionOptions = {
  automaticLayout: true,
  cursorStyle: 'line',
  extraEditorClassName: 'codeEditorMonaco',
  lineNumbersMinChars: 2,
  readOnly: false,
  roundedSelection: false,
  selectOnLineNumbers: true,
}

// CodeViewer wraps the code renderer.
export class CodeViewer extends React.Component<
  ICodeViewerProps,
  ICodeViewerState
> {
  constructor(props: ICodeViewerProps) {
    super(props)
    this.state = {}
  }

  // render draws the object with the desired renderer
  // TODO: properly multiplex available UI components
  // for now show a Prism editor with the JSON
  public render() {
    return (
      <MonacoEditor
        width="100%"
        height="100%"
        language={this.props.language || 'json'}
        options={{...editorOptions, readOnly: true}}
        value={this.props.data || ''}
      />
    )
  }
}
