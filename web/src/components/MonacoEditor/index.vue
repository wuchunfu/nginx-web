<template>
  <div class="container">
    <div class="code" ref="containerDomRef"></div>
  </div>
</template>

<script lang="ts">
import { onBeforeUnmount, onMounted, PropType, reactive, ref, shallowRef, toRefs, watch,inject } from 'vue'

import * as Monaco from 'monaco-editor'

// 注册 TypeScript 语言服务，提供代码补全、查找引用、重命名等语言能力
// import "monaco-typescript/release/esm/monaco.contribution";
// 注册 TypeScript 语言的语法解释器，提供语法高亮等，是 TypeScript 语言的基础能力
// import "monaco-languages/release/esm/typescript/typescript.contribution";
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import JsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import CssWorker from 'monaco-editor/esm/vs/language/css/css.worker?worker'
import HtmlWorker from 'monaco-editor/esm/vs/language/html/html.worker?worker'
import TsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import { IRange } from "monaco-editor";

const global: any = globalThis || window

// 构建 Monaco Editor 的运行环境，因为 WebWorker 必须使用独立文件的方式引入，所以 Monaco Editor 使
// 用了一个相对原始的方式引入其 worker ，注意需要使用 ?worker 来引入 Web Worker 。
// https://github.com/vitejs/vite/discussions/1791#discussioncomment-321046
global.MonacoEnvironment = {
  getWorker(_: string, label: string) {
    if (label === 'json') {
      return new JsonWorker()
    }
    if (label === 'css' || label === 'scss' || label === 'less') {
      return new CssWorker()
    }
    if (label === 'html' || label === 'handlebars' || label === 'razor') {
      return new HtmlWorker()
    }
    if (label === 'typescript' || label === 'javascript') {
      return new TsWorker()
    }
    return new EditorWorker()
  }
}

const languages = Monaco.languages.getLanguages()

let subscription: Monaco.IDisposable | undefined
let preventTriggerChangeEvent = false

export default {
  name: 'MonacoEditor',
  props: {
    isToolbar: {
      type: Boolean,
      default: true
    },
    code: {
      // 代码
      type: String as PropType<string>,
      required: true,
      default: '{ "code": "hello world!" }'
    },
    // selection: {
    //   type: String as PropType<string>,
    //   required: false,
    //   default: ''
    // },
    layout: {
      // 布局
      type: Object as PropType<Monaco.editor.IDimension>,
      required: true,
      default: () => ({})
    },
    width: {
      type: [String, Number],
      default: '100%'
    },
    height: {
      type: [String, Number],
      default: '100vh'
    },
    theme: {
      type: String,
      default: 'vs-dark'
    },
    language: {
      type: String,
      default: 'json'
    },
    languageModal: {
      type: Array,
      default: () => [
        'go',
        'java',
        'python',
        'rust',
        'kotlin',
        'shell',
        'sql',
        'json',
        'xml',
        'yaml',
        'mysql',
        'pgsql',
        'redis',
        'markdown',
        'typescript',
        'javascript',
        'html',
        'css',
        'scss',
        'less',
        'c',
        'cpp',
        'ini',
        'graphql',
        'scheme',
        'powershell',
        'dockerfile',
        'apex',
        'azcli',
        'bat',
        'clojure',
        'coffeescript',
        'csharp',
        'csp',
        'fsharp',
        'handlebars',
        'lua',
        'msdax',
        'objective-c',
        'pascal',
        'perl',
        'php',
        'plaintext',
        'postiats',
        'powerquery',
        'pug',
        'r',
        'razor',
        'redshift',
        'ruby',
        'sb',
        'sol',
        'st',
        'swift',
        'tcl',
        'vb',
      ],
    },
    options: {
      type: Object as PropType<Monaco.editor.IStandaloneEditorConstructionOptions>,
      default: () => ({
        isToolbar: true,
        wordWrap: true, //on,off
        foldingStrategy: "indentation", // 代码可分小段折叠
        automaticLayout: true, // 自适应布局
        overviewRulerBorder: false, // 不要滚动条的边框
        autoClosingBrackets: true,
        mouseWheelZoom: false,
        tabSize: 2, // tab 缩进长度
        fontSize: 14,
        minimap: {
          enabled: true, // 不要小地图
        },
        lineNumbers: true, //on,off
      })
    },
    vid: [String, Number],
    onChange: {
      type: Function as PropType<(value: string, event: Monaco.editor.IModelContentChangedEvent) => void>
    },
    title: {
      type: String as PropType<string>,
      default: ''
    },
    // 是否使用diff模式
    diffEditor: {
      type: Boolean,
      default: false
    },
    // 只有在diff模式下有效
    original: {
      type: String,
      default: "json"
    }
  },
  setup: (props: any, { emit }: any) => {
    // 需要一个shallowRef: 只监听value，不关心实际对象
    const editorRef = shallowRef<Monaco.editor.IStandaloneCodeEditor | null>(null)

    // 需要生成编辑器的Dom
    const containerDomRef = ref<HTMLElement | null>(null)

    const state = reactive({
      // 主要配置
      defaultOptions: {
        // 编辑器的初始值
        value: props.code,
        // 是否可编辑
        readOnly: false,
        // 编辑器主题：vs, hc-black, or vs-dark，更多选择详见官网
        // vs, hc-black, or vs-dark
        theme: 'vs-dark',
        // 代码生成语言
        language: 'json',
        // 当粘贴的时候自动进行一次格式化代码
        formatOnPaste: true,
        // tab缩进长度
        tabSize: 2,
        minimap: {
          // 不需要小的缩略图
          enabled: false
        },
        // 右侧不显示编辑器预览框
        roundedSelection: true,
        // 自动缩进
        autoIndent: true,
        // 字体
        fontFamily: '微软雅黑',
        // 默认的提示关掉
        quickSuggestions: true,
        // 代码提示延时
        quickSuggestionsDelay: 100,
        // 编辑器自适应布局，可能会影响性能
        // automaticLayout: true,
        overviewRulerBorder: false,
        // 滚动配置，溢出才滚动
        scrollBeyondLastLine: false,
      },
    });

    // 格式化代码
    const formatCode = () => {
      window.requestIdleCallback(() => {
          editorRef.value!.getAction('editor.action.formatDocument').run()
        }, {
          timeout: 600
        }
      )
    }

    onMounted(() => {
      // 生成编辑器配置
      const editorOptions = Object.assign(state.defaultOptions, props.options)

      // 组件初始化时创建一个MonacoEditor的实例
      // editorRef.value = Monaco.editor.create(containerDomRef.value!, {
      //   value: props.code, // 初始值
      //   theme: 'vs-dark', // vs, hc-black, or vs-dark
      //   language: 'json', // 代码生成语言
      //   formatOnPaste: true, // 当粘贴的时候自动进行一次格式化代码
      //   tabSize: 2, // tab缩进长度
      //   minimap: {
      //     enabled: false // 不需要小的缩略图
      //   },
      //   fontFamily: '微软雅黑', //字体
      //   // automaticLayout: true, //编辑器自适应布局，可能会影响性能
      //   overviewRulerBorder: false,
      //   autoIndent: true ,
      //   scrollBeyondLastLine: false, //滚动配置，溢出才滚动
      //   ...props.options
      // })
      editorRef.value = Monaco.editor.create(containerDomRef.value!, editorOptions)

      // props.onChange?.(editorRef.value!.getValue())

      // emit("onChange", editorRef.value!.getValue());
      //
      // const selection: any = editorRef.value?.getSelection();
      // const range = new Monaco.Range(
      //   selection.startLineNumber,
      //   selection.startColumn,
      //   selection.endLineNumber,
      //   selection.endColumn,
      // )
      //
      // let valueInRange = editorRef.value!.getModel()?.getValueInRange(range);
      //
      // console.log(333333)
      // console.log(selection)
      // console.log(valueInRange)
      // props.onSelection?.(valueInRange)

      // emit("selection", valueInRange);

      // 如果代码有变化，会在这里监听到，当受到外部数据改变时，不需要触发change事件
      subscription = editorRef.value.onDidChangeModelContent((event) => {
        if (!preventTriggerChangeEvent) {
          // getValue: 获取编辑器中的所有文本
          props.onChange?.(editorRef.value!.getValue(), event)
          emit("onChange", editorRef.value!.getValue());
          // const selection: any = editorRef.value?.getSelection();
          // let valueInRange = editorRef.value!.getModel()?.getValueInRange(selection);

          // emit("onSelection", valueInRange);
          // emit("selection", valueInRange);

          emit("onChange", editorRef.value!.getValue());

          // emit("onSelection", editorRef.value?.getSelection());


          // const selection: any = editorRef.value?.getSelection();
          // const range = new Monaco.Range(
          //   selection.startLineNumber,
          //   selection.startColumn,
          //   selection.endLineNumber,
          //   selection.endColumn,
          // )
          //
          // let valueInRange = editorRef.value!.getModel()?.getValueInRange(range);
          //
          // console.log(333333)
          // console.log(selection)
          // console.log(valueInRange)
          // props.onSelection?.(valueInRange)
        }
      })
      formatCode()
      editorRef.value.layout(props.layout)
    })

    onBeforeUnmount(() => {
      // 组件销毁时卸载编辑器
      if (subscription) {
        subscription.dispose()
      }
    })

    // 更新编辑器
    const refreshEditorHandler = () => {
      if (editorRef.value) {
        const editor = editorRef.value
        // 获取编辑器的textModel文本
        const model = editor.getModel()

        // 如果代码发生变化 这里需要更新一版
        if (model && props.code !== model.getValue()) {
          // 这是进行一次常规化的操作 文档原文：Push an "undo stop" in the undo-redo stack.
          editor.pushUndoStop()
          preventTriggerChangeEvent = true
          /**
           * @function 开始编辑编辑器, 文档原文：Push edit operations, basically editing the model. This is the preferred way of editing the model. The edit operations will land on the undo stack.
           * @param 编辑操作之前的光标状态。调用撤销或重做时，将返回此光标状态
           * @param 需要编辑的内容 range: 编辑的内容范围，这里选择的是全部范围
           * @param 在编辑器质性完成之后可以计算光标状态的一个回调参数
           */
          model.pushEditOperations([], [
              {
                range: model.getFullModelRange(),
                text: props.code
              }
            ], () => null
          )
        }

        editor.pushUndoStop()
        preventTriggerChangeEvent = false
        formatCode()
      }
    }

    const themeHandler = () => {
      if (editorRef.value) {
        Monaco.editor.setTheme(props.theme)
      }
    };

    const languageHandler = (val: any) => {
      if (editorRef.value) {
        Monaco.editor.setTheme(props.theme)
      }
    };

    watch(() => props.code, refreshEditorHandler, {
      immediate: true,
      deep: true
    });

    watch(() => props.theme, themeHandler, {
      immediate: true,
      deep: true
    });

    watch(() => props.language, () => {
      if (editorRef.value) {
        let langModel: any = editorRef.value.getModel();
        Monaco.editor.setModelLanguage(langModel, props.language);
      }
    }, {
      immediate: true,
      deep: true
    });

    return {
      ...toRefs(state),
      editorRef,
      containerDomRef
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  border: 1px solid #eee;
  display: flex;
  flex-direction: column;
  border-radius: 5px;
}

.title {
  background-color: #eee;
  padding: 10px 0 10px 20px;
}

.code {
  flex-grow: 1
}
</style>
