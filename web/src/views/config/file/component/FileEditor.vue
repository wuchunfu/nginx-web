<template>
  <div class="system-add-user-container">
    <el-card shadow="hover">
      <el-row :gutter="10" class="mb8">
        <el-col :span="1.5">
          <el-button
            type="success"
            size="small"
            @click="handleSave"
          >
            <el-icon>
              <ele-CirclePlus/>
            </el-icon>
            保存
          </el-button>
        </el-col>
      </el-row>

      <el-skeleton
        :loading="loading.fileEditor"
        animated
        :throttle="500"
      >
        <MonacoEditor
          :code="codes"
          :layout="opts.sty"
          @onChange="codeChange"
        />
      </el-skeleton>
    </el-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive, ref, toRefs, watch } from 'vue';
import MonacoEditor from '/@/components/MonacoEditor/index.vue'
import { useRoute } from "vue-router";
import { getDetail } from "/@/api/config";

interface FormDataState {
  loading: {
    fileEditor: boolean,
  },
  // 是否显示弹出层
  isShowDialog: boolean;
  // 弹出层标题
  title: string,
  codes: string,
  opts: object,
}

export default defineComponent({
  name: 'FileEditor',
  components: {
    MonacoEditor
  },
  setup(props, context) {

    const route = useRoute();

    const ruleFormRef = ref();

    const state = reactive<FormDataState>({
      loading: {
        fileEditor: false,
      },
      isShowDialog: false,
      title: '',
      codes: "",
      opts: {
        sty: { width: "100%", height: 300 },
        theme: "vs",
        language: "json",
      },
    });

    // 页面加载时
    onMounted(() => {
      handleDetail();
    });

    watch(() => route.query.fileName, () => {
      console.log(route.query.fileName)
    }, {
      immediate: true,
      deep: true
    });

    const codeChange = (code: string) => {
      state.codes = code;
    };

    const handleDetail = () => {
      const params = {
        parentPath: route.query.fileName
      }
      getDetail(params).then((res: any) => {
        console.log("detail: ", res)
        if (res.code === 200) {
          state.codes = res.data;
        }
      });
    }

    const handleSave = () => {

    }

    return {
      ...toRefs(state),
      ruleFormRef,
      codeChange,
      handleSave,
    };
  },
});
</script>
