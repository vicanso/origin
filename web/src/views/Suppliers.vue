<template>
  <div class="suppliers">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-goods" />
        供应商
      </div>
      <BaseFilter :fields="filterFields" v-if="inited" @filter="filter" />
      <el-table v-loading="processing" :data="suppliers" row-key="id" stripe>
        <el-table-column prop="name" key="name" label="名称" />
        <el-table-column
          prop="statusDesc"
          key="statusDesc"
          label="状态"
          width="80"
        />
        <el-table-column prop="address" key="address" label="地址" />
        <el-table-column prop="contact" key="contact" label="联系人" />
        <el-table-column prop="mobile" key="mobile" label="联系电话" />
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
        :total="supplierCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </el-card>
    <Supplier v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import Supplier from "@/components/products/Supplier.vue";

const supplierStatuses = [];
const filterFields = [
  {
    label: "状态：",
    key: "status",
    type: "select",
    options: supplierStatuses,
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
  name: "Suppliers",
  extends: BaseTable,
  components: {
    BaseFilter,
    Supplier
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
      processing: state => state.supplier.processing,
      suppliers: state => state.supplier.list.data || [],
      supplierCount: state => state.supplier.count
    })
  },
  methods: {
    ...mapActions(["listStatus", "listSupplier"]),
    async fetch() {
      const { query } = this;
      try {
        await this.listSupplier(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  },
  async beforeMount() {
    try {
      const { statuses } = await this.listStatus();
      supplierStatuses.length = 0;
      supplierStatuses.push({
        name: "所有",
        value: null
      });
      supplierStatuses.push(...statuses);
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
.suppliers
  margin: $mainMargin
.addBtn
  width: 100%
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
