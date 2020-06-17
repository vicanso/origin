<template>
  <div class="advertisements" v-loading="!inited">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-c-scale-to-original" />
        广告
      </div>
      <BaseFilter :fields="filterFields" v-if="inited" @filter="filter" />
      <el-table
        v-loading="processing"
        :data="advertisements"
        row-key="id"
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column
          prop="categoryDesc"
          key="categoryDesc"
          label="类别"
          width="80"
        />
        <el-table-column
          prop="statusDesc"
          key="statusDesc"
          label="状态"
          width="80"
        />
        <el-table-column
          prop="rank"
          key="rank"
          label="排序"
          width="60"
          sortable
        />
        <el-table-column prop="link" key="link" label="链接" width="200" />
        <el-table-column label="图片" width="80">
          <template slot-scope="scope">
            <img class="pic" :src="scope.row.pic" />
          </template>
        </el-table-column>
        <el-table-column
          prop="startedAtDesc"
          key="startedAtDesc"
          label="生效时间"
          width="180"
        />
        <el-table-column
          prop="endedAtDesc"
          key="endedAtDesc"
          label="生效时间"
          width="180"
        />
        <el-table-column
          prop="updatedAtDesc"
          key="updatedAtDesc"
          label="更新时间"
          width="180"
        />
        <el-table-column prop="summary" key="summary" label="简介" />
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
        :total="advertisementCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </el-card>
    <Advertisement v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import Advertisement from "@/components/Advertisement.vue";

const advertisementStatuses = [];

const filterFields = [
  {
    label: "状态：",
    key: "status",
    type: "select",
    options: advertisementStatuses,
    span: 8
  },
  {
    label: "关键字：",
    key: "keyword",
    placeholder: "请输入关键字",
    clearable: true,
    span: 8
  },
  {
    label: "",
    type: "filter",
    span: 8,
    labelWidth: "0px"
  }
];

export default {
  name: "Advertisements",
  extends: BaseTable,
  components: {
    BaseFilter,
    Advertisement
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
      processing: state => state.advertisement.processing,
      advertisementCount: state => state.advertisement.list.count,
      advertisements: state => state.advertisement.list.data || []
    })
  },
  methods: {
    ...mapActions(["listAdvertisement", "listStatus"]),
    async fetch() {
      const { query } = this;
      try {
        await this.listAdvertisement(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    try {
      const { statuses } = await this.listStatus();
      advertisementStatuses.length = 0;
      advertisementStatuses.push({
        name: "所有",
        value: null
      });
      advertisementStatuses.push(...statuses);
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
.advertisements
  margin: $mainMargin
  i
    margin-right: 5px
.pic
  max-height: 60px
.addBtn
  width: 100%
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
