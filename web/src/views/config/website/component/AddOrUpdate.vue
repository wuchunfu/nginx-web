<template>
  <div class="system-add-user-container">
    <el-dialog
      :title="formData.title"
      v-model="formData.isShowDialog"
      width="40%"
      destroy-on-close
      :before-close="handleClose"
    >
      <el-form
        ref="ruleFormRef"
        :rules="formData.rules"
        :model="formData.ruleForm"
        status-icon
        size="default"
        label-width="110px"
      >
        <el-form-item
          label="配置名称"
          prop="fileName"
          style="width: 90%;"
        >
          <el-input
            type="text"
            v-model="formData.ruleForm.fileName"
            placeholder="请输入配置名称"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item
          label="网站域名"
          prop="serverName"
          style="width: 90%;"
        >
          <el-input
            type="text"
            v-model="formData.ruleForm.serverName"
            placeholder="请输入网站域名"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item
          label="网站根目录"
          prop="rootDirectory"
          style="width: 90%;"
        >
          <el-input
            type="text"
            v-model="formData.ruleForm.rootDirectory"
            placeholder="请输入网站根目录"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item
          label="网站首页"
          prop="homePage"
          style="width: 90%;"
        >
          <el-input
            type="text"
            v-model="formData.ruleForm.homePage"
            placeholder="请输入网站首页"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item
          label="http 监听端口"
          prop="httpPort"
          style="width: 90%;"
        >
          <el-input-number
            v-model="formData.ruleForm.httpPort"
            :min="1"
            :max="65535"
            controls-position="right"
            @change="handleHttpChange"
          />
        </el-form-item>
        <el-form-item
          label="启用 TLS"
          prop="supportSsl"
          style="width: 90%;"
        >
          <el-switch
            v-model="formData.ruleForm.supportSsl"
            active-color="#13ce66"
            inactive-color="#ff4949"
            :active-value="1"
            :inactive-value="0"
            @change="handleSslChange"
          >
          </el-switch>
        </el-form-item>
        <el-form-item
          label="https 监听端口"
          prop="httpsPort"
          v-if="formData.ruleForm.supportSsl === 1"
          style="width: 90%;"
        >
          <el-input-number
            v-model="formData.ruleForm.httpsPort"
            :min="1"
            :max="65535"
            controls-position="right"
            @change="handleHttpsChange"
          />
        </el-form-item>
        <el-form-item
          label="TLS 证书路径"
          prop="sslCertificate"
          v-if="formData.ruleForm.supportSsl === 1"
          style="width: 90%;"
        >
          <el-input
            type="text"
            v-model="formData.ruleForm.sslCertificate"
            placeholder="请输入TLS 证书路径"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item
          label="私钥路径"
          prop="sslCertificateKey"
          v-if="formData.ruleForm.supportSsl === 1"
          style="width: 90%;"
        >
          <el-input
            type="text"
            v-model="formData.ruleForm.sslCertificateKey"
            placeholder="请输入私钥路径"
            clearable
          >
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button
            :loading="formData.loading.summitForm"
            type="primary"
            @click="onSubmit"
            size="default"
          >
            确定
          </el-button>
					<el-button
            @click="onCancel"
            size="default"
          >
            取消
          </el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive, ref, toRefs } from 'vue';
import { getDetailById, save, update } from "/@/api/website";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";

// 定义接口来定义对象的类型
interface FormDataRow {
  websiteId: number;
  fileName: string;
  serverName: string;
  rootDirectory: string;
  homePage: string;
  httpPort: number;
  supportSsl: number;
  httpsPort: number;
  sslCertificate: string;
  sslCertificateKey: string;
  status: number;
  createTime: string;
  updateTime: string;
}

interface FormDataState {
  formData: {
    loading: {
      summitForm: boolean,
    },
    // 是否显示弹出层
    isShowDialog: boolean;
    // 弹出层标题
    title: string,
    // 表单
    ruleForm: FormDataRow;
    // 表单校验
    rules: object;
  }
}

export default defineComponent({
  name: 'AddOrUpdateWebsite',
  setup(props, context) {
    const ruleFormRef = ref();

    const state = reactive<FormDataState>({
      formData: {
        loading: {
          summitForm: false,
        },
        isShowDialog: false,
        title: '',
        ruleForm: {} as FormDataRow,
        rules: {
          fileName: [{ required: true, message: "Please input config name", trigger: "blur" }],
          serverName: [{ required: true, message: "Please input server name", trigger: "blur" }],
          rootDirectory: [{ required: true, message: "Please input root directory", trigger: "blur" }],
          homePage: [{ required: true, message: "Please input home page", trigger: "blur" }],
          httpPort: [{ required: true, message: "Please input http port", trigger: "blur" }],
        },
      }
    });

    // 页面加载时
    onMounted(() => {
    });

    // 表单重置
    const reset = () => {
      state.formData.ruleForm = {
        websiteId: -1,
        fileName: '',
        serverName: '',
        rootDirectory: '',
        homePage: '',
        httpPort: 80,
        supportSsl: 0,
        httpsPort: 443,
        sslCertificate: '',
        sslCertificateKey: '',
        status: 0,
        createTime: '',
        updateTime: '',
      };
      ruleFormRef.value?.resetFields()
    };

    // 打开弹窗
    const openAddDialog = () => {
      reset();
      state.formData.isShowDialog = true;
      state.formData.title = "添加";
    };

    const openUpdateDialog = (fileName: string) => {
      reset();
      getDetailById(fileName).then((res: any) => {
        console.log("detail: ", res)
        if (res.code === 200) {
          state.formData.ruleForm = res.data;
          state.formData.isShowDialog = true;
          state.formData.title = "修改";
        }
      }).catch((res) => {
        console.log(res)
        ElNotification({
          type: 'error',
          showClose: true,
          duration: 3000,
          title: "详情",
          message: "获取详情失败"
        });
      });
    };

    // 关闭弹窗
    const closeDialog = () => {
      reset();
      state.formData.isShowDialog = false;
    };

    // 新增
    const onSubmit = () => {
      state.formData.loading.summitForm = true;
      ruleFormRef.value?.validate((valid: boolean) => {
        if (valid) {
          if (state.formData.ruleForm.websiteId !== -1) {
            update(state.formData.ruleForm).then((res: any) => {
              console.log("update: ", res)
              if (res.code === 200) {
                ElNotification({
                  type: 'success',
                  showClose: true,
                  duration: 3000,
                  title: "更新",
                  message: "更新成功"
                });
                context.emit("onRefreshDataList")
                state.formData.loading.summitForm = false;
              } else {
                ElNotification({
                  type: 'error',
                  showClose: true,
                  duration: 3000,
                  title: "更新",
                  message: "更新失败"
                });
                state.formData.loading.summitForm = false;
              }
              closeDialog();
            }).catch((res: any) => {
              console.log(res)
              ElNotification({
                type: 'error',
                showClose: true,
                duration: 3000,
                title: "更新",
                message: "更新失败"
              });
              state.formData.loading.summitForm = false;
            });
          } else {
            save(state.formData.ruleForm).then((res: any) => {
              console.log("save: ", res)
              if (res.code === 200) {
                ElNotification({
                  type: 'success',
                  showClose: true,
                  duration: 3000,
                  title: "保存",
                  message: "保存成功"
                });
                context.emit("onRefreshDataList")
                state.formData.loading.summitForm = false;
              } else {
                ElNotification({
                  type: 'success',
                  showClose: true,
                  duration: 3000,
                  title: "保存",
                  message: "保存失败"
                });
                state.formData.loading.summitForm = false;
              }
              closeDialog();
            }).catch((res: any) => {
              console.log(res)
              ElNotification({
                type: 'error',
                showClose: true,
                duration: 3000,
                title: "保存",
                message: "保存失败"
              });
              state.formData.loading.summitForm = false;
            });
          }
        } else {
          console.log('error submit!')
          state.formData.loading.summitForm = false;
          return false
        }
      });
    };

    // 取消
    const onCancel = () => {
      closeDialog();
    };

    const handleClose = (done: () => void) => {
      ElMessageBox.confirm('Are you sure to close this dialog?').then(() => {
        done()
      }).catch((res) => {
        console.log(res)
        ElMessage({
          type: 'warning',
          message: 'Cancel close this dialog',
          showClose: true,
        })
      })
    }

    const handleHttpChange = (value: number) => {
      state.formData.ruleForm.httpPort = value;
      console.log(value)
    }

    const handleSslChange = (value: number) => {
      console.log(value)
      console.log(state.formData.ruleForm)
      state.formData.ruleForm.httpsPort = 443;
      state.formData.ruleForm.sslCertificate = '';
      state.formData.ruleForm.sslCertificateKey = '';
    }

    const handleHttpsChange = (value: number) => {
      state.formData.ruleForm.httpsPort = value;
      console.log(value)
    }

    return {
      ...toRefs(state),
      ruleFormRef,
      openAddDialog,
      openUpdateDialog,
      closeDialog,
      onCancel,
      onSubmit,
      handleClose,
      handleHttpChange,
      handleSslChange,
      handleHttpsChange,
    };
  },
});
</script>
