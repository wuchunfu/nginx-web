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
            type="warning"
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
        <el-col :span="1.5">
          <el-button
            type="success"
            size="small"
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
          @cell-click="handleListFolder"
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
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive, ref, toRefs } from 'vue';
import { ElMessage, ElNotification } from 'element-plus';
import { getList, getListFolder } from "/@/api/config";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";

// 定义接口来定义对象的类型
interface TableDataRow {
  basePath: string;
  parentPath: string;
  filePath: string;
  fileName: string;
  isFile: boolean;
  fileType: string;
  fileSize: string;
  suffixName: string;
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
  name: 'userList',
  setup() {
    const { t } = useI18n();
    const router = useRouter();

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
        if (row.filePath === "") {
          state.tableData.fileParentPath = row.fileName;
        } else {
          state.tableData.fileParentPath = row.filePath;
        }
        state.tableData.folderDisabled = false;
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
          let data = res.data[0];
          if (data.parentPath === "") {
            state.tableData.folderDisabled = true
          }
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

    // 多选框选中数据
    const handleSelectionChange = (selection: any) => {
      state.tableData.fileNames = selection.map((item: any) => item.fileName)
      state.tableData.single = selection.length != 1
      state.tableData.multiple = !selection.length
    };

    // 打开修改用户弹窗
    const handleUpdate = (row: TableDataRow) => {
      let fileName = "";
      if (state.tableData.fileNames.length !== 0) {
        fileName = state.tableData.fileNames;
      } else {
        if (row.filePath === "") {
          fileName = row.fileName;
        } else {
          fileName = row.filePath;
        }
      }

      router.push({
        name: "fileEditor",
        query: {
          fileName: fileName
        }
      })
    };

    return {
      ...toRefs(state),
      queryFormRef,
      handleSelectionChange,
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
