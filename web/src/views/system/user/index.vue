<template>
  <div class="system-user-container">
    <el-card shadow="hover">
      <el-form
        :model="tableData.queryParams"
        ref="queryFormRef"
        :inline="true"
        v-show="tableData.showSearch"
        label-width="68px"
      >
        <el-form-item label="账户名称" prop="username">
          <el-input
            v-model="tableData.queryParams.username"
            placeholder="请输入账户名称"
            clearable
            size="default"
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="default"
            @click="handleQuery"
            :loading="tableData.tableLoading"
            class="ml10"
          >
            <el-icon>
              <ele-Search/>
            </el-icon>
            搜索
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button
            size="default"
            @click="resetQuery"
            class="ml10"
          >
            <el-icon>
              <ele-Refresh/>
            </el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>

      <el-row :gutter="10" class="mb8">
        <el-col :span="1.5">
          <el-button
            type="primary"
            plain
            size="default"
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
            plain
            size="default"
            :disabled="tableData.single"
            @click="handleUpdate"
          >
            <el-icon>
              <ele-Edit/>
            </el-icon>
            修改
          </el-button>
        </el-col>
        <el-col :span="1.5">
          <el-button
            type="danger"
            plain
            size="default"
            :disabled="tableData.multiple"
            @click="handleDelete"
          >
            <el-icon>
              <ele-Delete/>
            </el-icon>
            批量删除
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
            prop="username"
            label="账号名称"
            align="center"
            width="100"
            fixed
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            prop="email"
            label="邮箱"
            align="center"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            prop="isAdmin"
            label="角色"
            align="center"
          >
            <template v-slot="scope">
              <span>{{ roleConvert(scope.row.isAdmin) }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="status"
            label="状态"
            align="center"
            width="90"
          >
            <template v-slot="scope">
            <span>
              <el-tag
                :type="scope.row.status === 1 ? 'success' : 'danger'"
                effect="plain"
              >
                {{ statusConvert(scope.row.status) }}
              </el-tag>
            </span>
            </template>
          </el-table-column>
          <el-table-column
            prop="createTime"
            label="创建时间"
            show-overflow-tooltip
            sortable
            align="center"
            width="170"
          >
            <template v-slot="scope">
              <span>{{ scope.row.createTime }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="updateTime"
            label="更新时间"
            show-overflow-tooltip
            sortable
            align="center"
            width="170"
          >
            <template v-slot="scope">
              <span>{{ scope.row.updateTime }}</span>
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
                :disabled="scope.row.usernames === 'admin'"
                size="small"
                type="text"
                @click="handleUpdate(scope.row)"
              >
                修改
              </el-button>
              <el-button
                :disabled="scope.row.usernames === 'admin'"
                size="small"
                type="text"
                @click="handleDelete(scope.row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-skeleton>

      <el-pagination
        style="float: right;"
        class="mt15 mb15"
        small
        background
        :pager-count="5"
        :page-sizes="[5, 10, 20, 50, 100, 200, 500, 1000]"
        :total="tableData.total"
        :hide-on-single-page="false"
        layout="total, sizes, prev, pager, next, jumper"
        v-model:current-page="tableData.queryParams.page"
        v-model:page-size="tableData.queryParams.pageSize"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      >
      </el-pagination>
    </el-card>
    <AddOrUpdate ref="addOrUpdateUserRef" @onRefreshDataList="getDataList"/>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive, ref, toRefs } from 'vue';
import { ElMessageBox, ElNotification } from 'element-plus';
import AddOrUpdate from '/@/views/system/user/component/AddOrUpdate.vue';
import { deleteById, getList } from "/@/api/user";
import { useI18n } from "vue-i18n";

// 定义接口来定义对象的类型
interface TableDataRow {
  userId: number;
  username: string;
  password: string;
  email: string;
  isAdmin: number;
  status: number;
  createTime: string;
  updateTime: string;
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
    // 查询参数
    queryParams: {
      page: number,
      pageSize: number,
      username: string,
    },
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
        queryParams: {
          page: 1,
          pageSize: 10,
          username: '',
        },
      },
    });

    // 页面加载时
    onMounted(() => {
      getDataList();
    });

    const roleConvert = (state: number) => {
      switch (state) {
        case 1:
          return "管理员";
        case 0:
          return "普通用户";
        default:
          return "普通用户";
      }
    }

    const statusConvert = (state: number) => {
      switch (state) {
        case 1:
          return "启用";
        case 0:
          return "禁用";
        default:
          return "禁用";
      }
    }

    // 表格数据
    const getDataList = () => {
      state.tableData.tableLoading = true;
      const params = {
        page: state.tableData.queryParams.page,
        pageSize: state.tableData.queryParams.pageSize,
        username: state.tableData.queryParams.username,
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

    // 搜索按钮操作
    const handleQuery = () => {
      state.tableData.queryParams.page = 1;
      getDataList();
    };

    // 重置按钮操作
    const resetQuery = () => {
      queryFormRef.value?.resetFields()
      handleQuery();
    };

    // 多选框选中数据
    const handleSelectionChange = (selection: any) => {
      state.tableData.ids = selection.map((item: any) => item.userId)
      state.tableData.single = selection.length != 1
      state.tableData.multiple = !selection.length
    };

    // 打开新增用户弹窗
    const handleAdd = () => {
      addOrUpdateUserRef.value.openAddDialog();
    };

    // 打开修改用户弹窗
    const handleUpdate = (row: TableDataRow) => {
      const rowId = state.tableData.ids.length !== 0 ? state.tableData.ids : row.userId;
      addOrUpdateUserRef.value.openUpdateDialog(rowId);
    };

    // 删除按钮操作
    const handleDelete = (row: TableDataRow) => {
      const ids = state.tableData.ids.length !== 0 ? state.tableData.ids : [row.userId];
      ElMessageBox.confirm(`是否确认删除数据项?`, '提示', {
        type: "warning",
        confirmButtonText: "确定",
        cancelButtonText: "取消"
      }).then(() => {
        const params = {
          userIds: ids
        }
        deleteById(params).then((res: any) => {
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

    // 每页页数
    const handleSizeChange = (val: number) => {
      console.log(`每页 ${ val } 条`)
      state.tableData.queryParams.page = 1;
      state.tableData.queryParams.pageSize = val;
      getDataList();
    };

    // 当前页
    const handleCurrentChange = (val: number) => {
      console.log(`当前页: ${ val }`)
      state.tableData.queryParams.page = val;
      getDataList();
    };

    return {
      ...toRefs(state),
      handleQuery,
      resetQuery,
      handleSelectionChange,
      addOrUpdateUserRef,
      queryFormRef,
      handleAdd,
      handleUpdate,
      handleDelete,
      handleSizeChange,
      handleCurrentChange,
      roleConvert,
      statusConvert,
      getDataList,
    };
  },
});
</script>

<style scoped lang="scss">

</style>
