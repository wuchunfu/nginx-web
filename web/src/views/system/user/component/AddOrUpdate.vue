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
        label-width="90px"
      >
        <el-form-item label="账号名称" prop="username" style="width: 90%;">
          <el-input
            type="text"
            v-model="formData.ruleForm.username"
            placeholder="请输入账号名称"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item label="账户密码" prop="password" style="width: 90%;">
          <el-input
            type="password"
            v-model="formData.ruleForm.password"
            autocomplete="off"
            placeholder="请输入账号密码"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword" style="width: 90%;">
          <el-input
            type="password"
            v-model="formData.ruleForm.confirmPassword"
            autocomplete="off"
            placeholder="请输入确认密码"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email" style="width: 90%;">
          <el-input
            type="text"
            v-model="formData.ruleForm.email"
            placeholder="请输入邮箱"
            clearable
          >
          </el-input>
        </el-form-item>
        <el-form-item label="角色" prop="isAdmin" style="width: 90%;">
          <el-select
            v-model="formData.ruleForm.isAdmin"
            clearable
            filterable
            placeholder="请选择"
            style="width: 100%;"
          >
            <el-option
              v-for="item in formData.roleOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status" style="width: 90%;">
          <el-switch
            v-model="formData.ruleForm.status"
            active-color="#13ce66"
            inactive-color="#ff4949"
            :active-value="1"
            :inactive-value="0"
          >
          </el-switch>
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
import { getDetailById, save, update } from "/@/api/user";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";

// 定义接口来定义对象的类型
interface FormDataRow {
  userId: number;
  username: string;
  password: string;
  confirmPassword: string;
  email: string;
  isAdmin: number;
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
    roleOptions: Array<object>;
  }
}

export default defineComponent({
  name: 'AddOrUpdateUser',
  setup(props, context) {
    const validatePassword = (rule: any, value: any, callback: any) => {
      if (value === '') {
        callback(new Error('Please input the password'))
      } else {
        if (state.formData.ruleForm.confirmPassword !== '') {
          if (!ruleFormRef.value) {
            return
          }
          ruleFormRef.value.validateField('confirmPassword', () => null)
        }
        callback()
      }
    }

    const validateConfirmPassword = (rule: any, value: any, callback: any) => {
      if (value === '') {
        callback(new Error('Please input the password again'))
      } else if (value !== state.formData.ruleForm.password) {
        callback(new Error("Two inputs don't match!"))
      } else {
        callback()
      }
    }

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
          username: [{ required: true, message: "Please input username", trigger: "blur" }],
          password: [{ required: true, validator: validatePassword, trigger: "blur" }],
          confirmPassword: [{ required: true, validator: validateConfirmPassword, trigger: "blur" }],
          email: [
            { required: true, message: "Please input email address", trigger: "blur" },
            { type: 'email', message: 'Please input correct email address', trigger: ['blur', 'change'], },
          ],
        },
        roleOptions: [
          {
            label: '普通用户',
            value: 0
          },
          {
            label: '管理员',
            value: 1
          }
        ]
      }
    });

    // 页面加载时
    onMounted(() => {
    });

    // 表单重置
    const reset = () => {
      state.formData.ruleForm = {
        userId: -1,
        username: '',
        password: '',
        confirmPassword: '',
        email: '',
        isAdmin: 0,
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

    const openUpdateDialog = (rowId: number) => {
      reset();
      getDetailById(rowId).then((res: any) => {
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
          if (state.formData.ruleForm.userId != -1) {
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

    return {
      ...toRefs(state),
      ruleFormRef,
      openAddDialog,
      openUpdateDialog,
      closeDialog,
      onCancel,
      onSubmit,
      handleClose,
    };
  },
});
</script>
