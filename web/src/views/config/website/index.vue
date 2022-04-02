<template>
  <div class="system-user-container">
    <el-card shadow="hover">
      <el-row :gutter="10" class="mb8">
        <el-col :span="1.5">
          <el-button
            type="primary"
            size="default"
            plain
            @click="handleAdd"
          >
            <el-icon>
              <ele-FolderAdd/>
            </el-icon>
            新增
          </el-button>
        </el-col>
        <el-col :span="1.5">
          <el-button
            type="success"
            size="default"
            plain
            :disabled="tableData.single"
            @click="handleUpdate"
          >
            <el-icon>
              <ele-Edit/>
            </el-icon>
            修改
          </el-button>
        </el-col>
      </el-row>

      <el-skeleton
        :loading="tableData.tableLoading"
        animated
        :throttle="500"
      >
        <el-table
          v-loading="tableData.tableLoading"
          :data="tableData.tableList"
          border
          fit
          tooltip-effect="dark"
          highlight-current-row
          element-loading-text="Loading"
          @selection-change="handleSelectionChange"
          style="width: 100%"
        >
          <el-table-column
            type="selection"
            align="center"
            width="55"
          />
          <el-table-column
            type="index"
            label="序号"
            align="center"
            width="65"
            fixed
          >
          </el-table-column>
          <el-table-column
            prop="fileName"
            label="名称"
            align="center"
            fixed
            show-overflow-tooltip
          >
            <template #default="scope">
              <el-icon color="#409EFF" size="20">
                <ele-Document/>
              </el-icon>
              <span style="margin-left: 5px">{{ scope.row.fileName }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="isEnabled"
            label="是否启用"
            align="center"
          >
            <template v-slot="scope">
              <span>
                <el-tag
                  :type="scope.row.isEnabled ? 'success' : 'danger'"
                  effect="plain"
                >
                  {{ statusConvert(scope.row.isEnabled) }}
                </el-tag>
              </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="fileSize"
            label="文件大小"
            align="center"
          >
          </el-table-column>
          <el-table-column
            prop="dateTime"
            label="修改时间"
            show-overflow-tooltip
            sortable
            align="center"
            width="170"
          >
            <template v-slot="scope">
              <span>{{ scope.row.dateTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            align="center"
            width="135"
            fixed="right"
          >
            <template #default="scope">
              <el-button
                size="small"
                type="text"
                @click="handleUpdate(scope.row)"
              >
                修改
              </el-button>
              <el-button
                size="small"
                type="text"
                v-if="!scope.row.isEnabled"
                @click="handleEnable(scope.row)"
              >
                启用
              </el-button>
              <el-button
                size="small"
                type="text"
                v-if="scope.row.isEnabled"
                @click="handleDisable(scope.row)"
              >
                禁用
              </el-button>
              <el-button
                size="small"
                type="text"
                :disabled="!scope.row.isEnabled"
                @click="handleDelete(scope.row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-skeleton>
    </el-card>

    <AddOrUpdate ref="addOrUpdateUserRef" @onRefreshDataList="getDataList"/>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive, ref, toRefs } from 'vue';
import { ElMessageBox, ElNotification } from 'element-plus';
import { deleteById, disable, enable, getList } from "/@/api/website";
import { useI18n } from "vue-i18n";
import AddOrUpdate from '/@/views/config/website/component/AddOrUpdate.vue';

// 定义接口来定义对象的类型
interface TableDataRow {
  fileName: string;
  isEnabled: boolean;
  fileSize: string;
  dateTime: string;
}

interface TableDataState {
  tableData: {
    // 遮罩层
    tableLoading: boolean;
    // 表格数据
    tableList: Array<TableDataRow>;
    // 总条数
    total: number;
    // 选中数组
    fileNames: any,
    // 非单个禁用
    single: boolean,
    // 非多个禁用
    multiple: boolean,
    // 显示搜索条件
    showSearch: boolean,
    fileParentPath: string
    folderDisabled: boolean
  };
}

export default defineComponent({
  name: 'websiteList',
  components: { AddOrUpdate },
  setup() {
    const { t } = useI18n();

    const addOrUpdateUserRef = ref();
    const queryFormRef = ref();

    const state = reactive<TableDataState>({
      tableData: {
        tableLoading: false,
        tableList: [],
        total: 0,
        fileNames: [],
        single: true,
        multiple: true,
        showSearch: true,
        fileParentPath: '',
        folderDisabled: true,
      },
    });

    // 页面加载时
    onMounted(() => {
      getDataList();
    });

    const statusConvert = (state: boolean) => {
      switch (state) {
        case true:
          return "启用";
        case false:
          return "禁用";
        default:
          return "禁用";
      }
    };

    // 表格数据
    const getDataList = () => {
      state.tableData.tableLoading = true;
      const params = {
        parentPath: state.tableData.fileParentPath
      }
      getList(params).then((res: any) => {
        console.log("getList: ", res)
        if (res.code === 200) {
          state.tableData.tableList = res.data;
          state.tableData.total = res.total;
          state.tableData.tableLoading = false;
        } else {
          ElNotification({
            type: 'error',
            showClose: true,
            duration: 3000,
            title: "列表",
            message: "获取数据失败"
          });
          state.tableData.tableLoading = false;
        }
      }).catch((res: any) => {
        console.log(res)
        ElNotification({
          type: 'error',
          showClose: true,
          duration: 3000,
          title: "列表",
          message: "获取数据失败"
        });
        state.tableData.tableLoading = false;
      });
    };

    // 多选框选中数据
    const handleSelectionChange = (selection: any) => {
      state.tableData.fileNames = selection.map((item: any) => item.fileName)
      state.tableData.single = selection.length != 1
      state.tableData.multiple = !selection.length
    };

    // 打开新增用户弹窗
    const handleAdd = () => {
      addOrUpdateUserRef.value.openAddDialog();
    };

    // 打开修改用户弹窗
    const handleUpdate = (row: TableDataRow) => {
      const fileName = state.tableData.fileNames.length !== 0 ? state.tableData.fileNames : row.fileName;
      addOrUpdateUserRef.value.openUpdateDialog(fileName);
    };

    const handleEnable = (row: TableDataRow) => {
      const fileName = row.fileName;
      enable(fileName).then((res: any) => {
        console.log("enable: ", res)
        if (res.code === 200) {
          getDataList();
          ElNotification({
            type: 'success',
            showClose: true,
            duration: 3000,
            title: "启用",
            message: "启用成功"
          });
        }
      }).catch((res) => {
        console.log(res)
        ElNotification({
          type: 'error',
          showClose: true,
          duration: 3000,
          title: "启用",
          message: "启用失败"
        });
      });
    }

    const handleDisable = (row: TableDataRow) => {
      const fileName = row.fileName;
      disable(fileName).then((res: any) => {
        console.log("disable: ", res)
        if (res.code === 200) {
          getDataList();
          ElNotification({
            type: 'success',
            showClose: true,
            duration: 3000,
            title: "禁用",
            message: "禁用成功"
          });
        }
      }).catch((res) => {
        console.log(res)
        ElNotification({
          type: 'error',
          showClose: true,
          duration: 3000,
          title: "禁用",
          message: "禁用失败"
        });
      });
    }

    // 删除按钮操作
    const handleDelete = (row: TableDataRow) => {
      const fileName = row.fileName;
      ElMessageBox.confirm(`是否确认删除数据项?`, '提示', {
        type: "warning",
        confirmButtonText: "确定",
        cancelButtonText: "取消"
      }).then(() => {
        deleteById(fileName).then((res: any) => {
          console.log("delete: ", res)
          if (res.code === 200) {
            getDataList();
            ElNotification({
              type: 'success',
              showClose: true,
              duration: 3000,
              title: "删除",
              message: "删除成功"
            });
          } else {
            ElNotification({
              type: 'error',
              showClose: true,
              duration: 3000,
              title: "删除",
              message: "删除失败"
            });
          }
        });
      })
    };

    return {
      ...toRefs(state),
      addOrUpdateUserRef,
      queryFormRef,
      handleSelectionChange,
      handleAdd,
      handleUpdate,
      handleEnable,
      handleDisable,
      handleDelete,
      statusConvert,
      getDataList,
    };
  },
});
</script>

<style scoped lang="scss">

</style>
