<template>
  <div class="system-add-user-container">
    <el-card shadow="hover">
      <el-affix :offset="50">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button
              type="success"
              size="small"
              @click="handleSave"
            >
              <el-icon>
                <ele-Edit/>
              </el-icon>
              保存
            </el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button
              type="warning"
              size="small"
              @click="handleCancel"
            >
              <el-icon>
                <ele-Back/>
              </el-icon>
              取消
            </el-button>
          </el-col>
        </el-row>
      </el-affix>

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
import { defineComponent, onMounted, reactive, ref, toRefs } from 'vue';
import MonacoEditor from '/@/components/MonacoEditor/index.vue'
import { useRoute, useRouter } from "vue-router";
import { getDetail, update } from "/@/api/config";
import { ElNotification } from "element-plus";

interface FormDataState {
  loading: {
    fileEditor: boolean,
  },
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
    const router = useRouter();

    const ruleFormRef = ref();

    const state = reactive<FormDataState>({
      loading: {
        fileEditor: false,
      },
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

    const codeChange = (code: string) => {
      state.codes = code;
    };

    const handleDetail = () => {
      state.loading.fileEditor = true;
      const params = {
        parentPath: route.query.fileName
      }
      getDetail(params).then((res: any) => {
        console.log("detail: ", res)
        if (res.code === 200) {
          state.codes = res.data;
        }
        state.loading.fileEditor = false;
      }).catch((res) => {
        console.log(res)
        ElNotification({
          type: 'error',
          showClose: true,
          duration: 3000,
          title: "详情",
          message: "获取详情失败"
        });
        state.loading.fileEditor = false;
      });
    }

    const handleSave = () => {
      state.loading.fileEditor = true
      const params = {
        parentPath: route.query.fileName,
        content: state.codes
      }
      update(params).then((res: any) => {
        console.log("update: ", res)
        if (res.code === 200) {
          state.codes = res.data;
          ElNotification({
            type: 'success',
            showClose: true,
            duration: 3000,
            title: "更新",
            message: "更新成功"
          });
        } else {
          ElNotification({
            type: 'error',
            showClose: true,
            duration: 3000,
            title: "更新",
            message: "更新失败"
          });
        }
        state.loading.fileEditor = false
      }).catch((res) => {
        console.log(res)
        ElNotification({
          type: 'error',
          showClose: true,
          duration: 3000,
          title: "更新",
          message: "更新失败"
        });
        state.loading.fileEditor = false;
      });
    };

    const handleCancel = () => {
      router.go(-1);
    };

    return {
      ...toRefs(state),
      ruleFormRef,
      codeChange,
      handleSave,
      handleCancel,
    };
  },
});
</script>
