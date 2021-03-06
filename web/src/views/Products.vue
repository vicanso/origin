<template>
  <div class="products">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-files" />
        产品
      </div>

      <el-table
        v-loading="processing"
        :data="products"
        row-key="id"
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="name" key="name" label="名称" width="180" />
        <el-table-column
          prop="price"
          key="price"
          label="单价"
          width="100"
          sortable
        />
        <el-table-column key="unit" label="单位" width="100">
          <template slot-scope="scope">
            {{ scope.row.specs + scope.row.unit }}
          </template>
        </el-table-column>
        <el-table-column
          prop="statusDesc"
          key="statusDesc"
          label="状态"
          width="80"
        />
        <el-table-column prop="brandDesc" key="brand" label="品牌" />
        <el-table-column prop="rank" key="rank" label="排序" sortable />
        <el-table-column
          prop="startedAtDesc"
          key="startedAtDesc"
          label="开始时间"
          width="180"
        />
        <el-table-column
          prop="endedAtDesc"
          key="endedAtDesc"
          label="结束时间"
          width="180"
        />
        <el-table-column label="分类" width="100">
          <template slot-scope="scope">
            <ul>
              <li v-for="category in scope.row.categoriesDesc" :key="category">
                {{ category }}
              </li>
            </ul>
          </template>
        </el-table-column>
        <el-table-column key="catalog" label="简介">
          <template slot-scope="scope">
            <el-tooltip placement="bottom">
              <pre slot="content">{{ scope.row.catalog }}</pre>
              <i class="el-icon-info" />
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="sn" key="sn" label="SN" />
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
        :total="productCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </el-card>
    <Product v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import Product from "@/components/products/Product.vue";

export default {
  name: "Products",
  extends: BaseTable,
  components: {
    Product
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
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
      processing: state => state.product.processing,
      productCount: state => state.product.list.count,
      products: state => state.product.list.data || []
    })
  },
  methods: {
    ...mapActions(["listProduct"]),
    async fetch() {
      const { query } = this;
      try {
        await this.listProduct(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.products
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
