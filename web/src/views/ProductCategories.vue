<template>
  <div class="productCategories">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-set-up" />
        产品分类
      </div>
      <BaseFilter :fields="filterFields" v-if="inited" @filter="filter" />
      <el-table
        v-loading="processing"
        :data="productCategories"
        row-key="id"
        stripe
      >
        <el-table-column prop="name" key="name" label="名称" />
        <el-table-column
          prop="statusDesc"
          key="statusDesc"
          label="状态"
          width="80"
        />
        <el-table-column prop="hot" key="hot" label="热度" width="80" />
        <el-table-column prop="level" key="level" label="级别" width="80" />
        <el-table-column key="belongs" label="所属分类">
          <template slot-scope="scope">
            <ul v-if="scope.row.belongsDesc">
              <li v-for="item in scope.row.belongsDesc" :key="item">
                {{ item }}
              </li>
            </ul>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column key="icon" label="图标">
          <template slot-scope="scope">
            <img v-if="scope.row.icon" :width="60" :src="scope.row.icon" />
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="updatedAtDesc"
          key="updatedAtDesc"
          label="更新时间"
          width="180"
        />
        <el-table-column fixed="right" label="操作" width="80">
          <template slot-scope="scope">
            <el-button
              class="op"
              type="text"
              size="small"
              @click="modify(scope.row)"
              >编辑</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        class="pagination"
        layout="prev, pager, next, sizes"
        :current-page="currentPage"
        :page-size="query.limit"
        :total="productCategoryCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </el-card>
    <ProductCategory v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import ProductCategory from "@/components/products/Category.vue";
import BaseFilter from "@/components/base/Filter.vue";

const categoryStatuses = [];
const filterFiedls = [
  {
    label: "状态：",
    key: "status",
    type: "select",
    options: categoryStatuses,
    span: 6
  },
  {
    label: "级别：",
    key: "level",
    type: "select",
    options: [
      {
        name: "所有",
        value: null
      },
      {
        name: "一级",
        value: 1
      },
      {
        name: "二级",
        value: 2
      },
      {
        name: "三级",
        value: 3
      }
    ],
    span: 6
  },
  {
    label: "关键字：",
    key: "keyword",
    placeholder: "请输入关键字",
    clearable: true,
    span: 6
  },
  {
    label: "",
    type: "filter",
    span: 6,
    labelWidth: "0px"
  }
];

export default {
  name: "ProductCategories",
  extends: BaseTable,
  components: {
    ProductCategory,
    BaseFilter
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      inited: false,
      filterFields: null,
      pageSizes,
      query: {
        offset: 0,
        limit: pageSizes[0],
        order: "-updatedAt"
      }
    };
  },
  computed: {
    ...mapState({
      processing: state => state.productCategory.processing,
      productCategoryCount: state => state.productCategory.list.count,
      productCategories: state => state.productCategory.list.data || []
    })
  },
  methods: {
    ...mapActions(["listProductCategory", "listStatus"]),
    async fetch() {
      const { query } = this;
      try {
        await this.listProductCategory(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    try {
      const { statuses } = await this.listStatus();
      categoryStatuses.length = 0;
      categoryStatuses.push({
        name: "所有",
        value: null
      });
      categoryStatuses.push(...statuses);
      this.filterFields = filterFiedls;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.inited = true;
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.productCategories
  margin: $mainMargin
.addBtn
  width: 100%
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
