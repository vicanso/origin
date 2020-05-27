<template>
  <div class="brands">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-goods" />
        产品品牌
      </div>
      <el-table v-loading="processing" :data="brands" row-key="id" stripe>
        <el-table-column prop="name" key="name" label="名称" width="120" />
        <el-table-column
          prop="statusDesc"
          key="statusDesc"
          label="状态"
          width="80"
        />
        <el-table-column prop="catalog" key="catalog" label="简介" />
        <el-table-column label="LOGO">
          <template slot-scope="scope">
            <img class="logo" :src="scope.row.logo" />
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
        :total="brandCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </el-card>

    <Brand v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import Brand from "@/components/products/Brand.vue";
import BaseTable from "@/components/BaseTable.vue";

export default {
  name: "Brands",
  extends: BaseTable,
  components: {
    Brand
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
      processing: state => state.brand.processing,
      brandCount: state => state.brand.list.count,
      brands: state => state.brand.list.data || []
    })
  },
  methods: {
    ...mapActions(["listBrand"]),
    async fetch() {
      const { query } = this;
      try {
        await this.listBrand(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.brands
  margin: $mainMargin
  i
    margin-right: 5px
.logo
  max-height: 60px
.submit, .selector, .addBtn
  width: 100%
.addBtn
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
