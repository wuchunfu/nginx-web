<template>
  <div class="system-user-container">
    <el-card shadow="hover">
      <el-row :gutter="10" class="mb8">
        <el-col :span="1.5">
          <el-button
            type="primary"
            size="small"
            :disabled="tableData.folderDisabled"
            @click="handleChangeFolder('all')"
          >
            <el-icon color="#FFFFFF" size="20">
              <ele-HomeFilled/>
            </el-icon>
            主页
          </el-button>
        </el-col>
        <el-col :span="1.5">
          <el-button
            type="primary"
            size="small"
            :disabled="tableData.folderDisabled"
            @click="handleChangeFolder('previous')"
          >
            <el-icon color="#FFFFFF" size="20">
              <ele-Back/>
            </el-icon>
            上一级
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
          @cell-click="handleListFolder"
          style="width: 100%"
        >
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
            align="left"
            fixed
            show-overflow-tooltip
          >
            <template #default="scope">
              <div v-if="scope.row.isFile">
                <el-icon color="#409EFF" size="20">
                  <ele-Document/>
                </el-icon>
                <span style="margin-left: 5px">{{ scope.row.fileName }}</span>
              </div>
              <div v-else>
                <el-icon color="#409EFF" size="20">
                  <ele-Folder/>
                </el-icon>
                <span style="margin-left: 5px;">{{ scope.row.fileName }}</span>
              </div>
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
            width="100"
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
import { ElMessage, ElNotification } from 'element-plus';
import AddOrUpdate from '/@/views/system/user/component/AddOrUpdate.vue';
import { getList, getListFolder } from "/@/api/config";
import { useI18n } from "vue-i18n";

// 定义接口来定义对象的类型
interface TableDataRow {
  userId: number;
  fileName: string;
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
    ids: any,
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
  name: 'userList',
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
        ids: [],
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

    const handleListFolder = (row: any, column: any, cell: any, event: any) => {
      if (!row.isFile) {
        state.tableData.folderDisabled = false;
        state.tableData.fileParentPath = row.filePath;
        getDataList();
      }
    };

    const handleChangeFolder = (flag: string) => {
      if ("all" === flag) {
        state.tableData.fileParentPath = "";
        state.tableData.folderDisabled = true;
        getDataList();
      } else if ("previous" === flag) {
        const params = {
          parentPath: state.tableData.fileParentPath
        }
        getListFolder(params).then((res: any) => {
          console.log("getListFolder: ", res);
          console.log("getListFolder: ", res.data[0]);
          // let data = JSON.parse(JSON.stringify(res.data[0]));
          let data = res.data[0];
          if (data.parentPath === data.basePath) {
            state.tableData.folderDisabled = true
          }
          // state.tableData.folderDisabled = true
          state.tableData.fileParentPath = data.parentPath
          state.tableData.tableList = res.data;
        }).catch((res: any) => {
          console.log(res);
          ElMessage({
            type: 'error',
            message: '文件夹切换失败',
          });
        });
      }
    };

    // 打开修改用户弹窗
    const handleUpdate = (row: TableDataRow) => {
      const rowId = state.tableData.ids.length !== 0 ? state.tableData.ids : row.userId;
      addOrUpdateUserRef.value.openUpdateDialog(rowId);
    };

    return {
      ...toRefs(state),
      addOrUpdateUserRef,
      queryFormRef,
      handleUpdate,
      getDataList,
      handleListFolder,
      handleChangeFolder,
    };
  },
});
</script>

<style scoped lang="scss">

</style>
