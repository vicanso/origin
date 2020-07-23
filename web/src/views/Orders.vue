<template>
  <div class="orders">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-files" />
        订单
      </div>
      <BaseFilter :fields="filterFields" v-if="inited" @filter="filter" />
      <el-table
        v-loading="processing"
        :data="orders"
        row-key="id"
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="sn" key="sn" label="订单编号" width="280" />
        <el-table-column
          prop="payAmount"
          key="payAmount"
          label="支付金额"
          width="100"
          sortable
        />
        <el-table-column
          prop="paySource"
          key="paySource"
          label="支付渠道"
          width="100"
        />
        <el-table-column
          prop="statusDesc"
          key="statusDesc"
          label="订单状态"
          width="90"
        />
        <el-table-column
          prop="delivererName"
          key="delivererName"
          label="送货员"
          width="100"
        />
        <el-table-column
          prop="receiverName"
          key="receiverName"
          label="收货人"
          width="120"
        />
        <el-table-column
          prop="receiverMobile"
          key="receiverMobile"
          label="收货手机"
          width="130"
        />
        <el-table-column
          prop="address"
          key="address"
          label="收货地址"
          width="300"
        />
        <el-table-column
          prop="createdAt"
          key="createdAt"
          label="创建时间"
          width="180"
          sortable
        />
        <el-table-column fixed="right" label="操作" width="80">
          <template slot-scope="scope">
            <el-button
              class="op"
              type="text"
              size="small"
              @click="modify(scope.row)"
              >详情</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        class="pagination"
        layout="prev, pager, next, sizes"
        :current-page="currentPage"
        :page-size="query.limit"
        :total="orderCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import { today, formatBegin, formatEnd } from "@/helpers/util";

const defaultDateRange = [today(), today()];
const orderStatuses = [];
const filterFields = [
  {
    label: "状态：",
    key: "statuses",
    type: "select",
    options: orderStatuses,
    multiple: true,
    span: 10
  },
  {
    label: "创建时间：",
    key: "dateRange",
    type: "dateRange",
    placeholder: ["开始日期", "结束日期"],
    defaultValue: defaultDateRange,
    span: 10
  },
  {
    label: "",
    type: "filter",
    span: 4,
    labelWidth: "0px"
  }
];

export default {
  name: "Orders",
  extends: BaseTable,
  components: {
    BaseFilter
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      pageSizes,
      inited: false,
      filterFields: null,
      query: {
        dateRange: defaultDateRange,
        offset: 0,
        limit: pageSizes[0],
        order: "-createdAt"
      }
    };
  },
  computed: {
    ...mapState({
      processing: state => state.order.processing,
      orderCount: state => state.order.list.count,
      orders: state => state.order.list.data || []
    })
  },
  methods: {
    ...mapActions(["listOrder", "listOrderStatus", "getOrderBySN"]),
    async fetch() {
      const query = Object.assign({}, this.query);
      const value = query.dateRange;
      if (value) {
        query.begin = formatBegin(value[0]);
        query.end = formatEnd(value[1]);
      } else {
        delete query.begin;
        delete query.end;
      }
      const { statuses } = query;
      delete query.statuses;
      if (statuses) {
        query.statuses = statuses.join(",");
      }
      delete query.dateRange;
      try {
        await this.listOrder(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    modify(item) {
      this.$router.push({
        query: {
          mode: CONFIG_EDITE_MODE,
          sn: item.sn
        }
      });
    }
  },
  async beforeMount() {
    try {
      const { statuses } = await this.listOrderStatus();
      orderStatuses.length = 0;
      orderStatuses.push(...statuses);
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
.orders
  margin: $mainMargin
  i
    margin-right: 5px
.submit, .selector, .addBtn
  width: 100%
.addBtn
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
