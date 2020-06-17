<template>
  <div class="regions">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-position" />
        地区
      </div>
      <BaseFilter :fields="filterFields" v-if="inited" @filter="filter" />
      <el-table
        v-loading="processing"
        :data="regions"
        row-key="id"
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="name" key="name" label="名称" />
        <el-table-column prop="code" key="code" label="代码" sortable />
        <el-table-column prop="statusDesc" key="statusDesc" label="状态" />
        <el-table-column
          prop="priority"
          key="priority"
          label="优先级"
          sortable
        />
        <el-table-column
          prop="updatedAtDesc"
          key="updatedAtDesc"
          label="更新时间"
          width="180"
        />
        <el-table-column label="操作" width="120">
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
        :total="regionCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>
    <RegionEditor v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import RegionEditor from "@/components/region/Editor.vue";
import BaseFilter from "@/components/base/Filter.vue";

const regionCategories = [];
const regionStatuses = [];
const filterFields = [
  {
    label: "分类：",
    key: "category",
    type: "select",
    options: regionCategories,
    span: 6
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    options: regionStatuses,
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
  name: "Regions",
  extends: BaseTable,
  components: {
    RegionEditor,
    BaseFilter
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      inited: false,
      fields: null,
      filterFields: null,
      pageSizes,
      query: {
        offset: 0,
        limit: pageSizes[0],
        order: "-updatedAt"
      }
    };
  },
  computed: mapState({
    processing: state => state.region.processing,
    regionCount: state => state.region.list.count,
    regions: state => state.region.list.data || []
  }),
  methods: {
    ...mapActions(["listRegion", "listStatus", "listRegionCategory"]),
    async fetch() {
      const { query } = this;
      try {
        await this.listRegion({
          params: query
        });
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    try {
      const { categories } = await this.listRegionCategory();
      regionCategories.length = 0;
      regionCategories.push({
        name: "所有",
        value: null
      });
      regionCategories.push(...categories);

      const { statuses } = await this.listStatus();
      regionStatuses.length = 0;
      regionStatuses.push({
        name: "所有",
        value: null
      });
      regionStatuses.push(...statuses);

      this.filterFields = filterFields;
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
.regions
  margin: $mainMargin
  i
    margin-right: 5px
.addBtn
  width: 100%
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
