<template>
  <div class="regions">
    <el-card v-if="!editMode">
      <div slot="header">
        <i class="el-icon-position" />
        地区
      </div>
      <el-form label-width="80px">
        <el-row :gutter="15">
          <el-col :span="6">
            <el-form-item label="分类：">
              <el-select
                class="selector"
                v-model="query.category"
                placeholder="请选择地区分类"
              >
                <el-option key="all-category" label="所有" value="" />
                <el-option
                  v-for="item in categories"
                  :key="item.value"
                  :label="item.name"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="状态：">
              <el-select
                class="selector"
                v-model="query.status"
                placeholder="请选择状态"
              >
                <el-option key="all-status" label="所有" value="" />
                <el-option
                  v-for="item in statuses"
                  :key="item.value"
                  :label="item.name"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="关键字：">
              <el-input
                clearable
                v-model="query.keyword"
                placeholder="请输入搜索关键字"
                @keyup.enter.native="handleSearch"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label-width="8px">
              <el-button
                :loading="processing"
                icon="el-icon-search"
                class="submit"
                type="primary"
                @click="handleSearch"
                >查询</el-button
              >
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <el-table v-loading="processing" :data="regions" row-key="id" stripe>
        <el-table-column prop="name" key="name" label="名称" />
        <el-table-column prop="code" key="code" label="代码" />
        <el-table-column prop="statusDesc" key="statusDesc" label="状态" />
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
import BaseTable from "@/components/BaseTable.vue";
import RegionEditor from "@/components/region/Editor.vue";

export default {
  name: "Regions",
  extends: BaseTable,
  components: {
    RegionEditor
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      pageSizes,
      currentRegion: null,
      originRegion: null,
      query: {
        offset: 0,
        limit: pageSizes[0],
        order: "-updatedAt",
        category: "",
        status: "",
        keyword: ""
      }
    };
  },
  computed: mapState({
    processing: state => state.region.processing,
    regionCount: state => state.region.list.count,
    regions: state => state.region.list.data || [],
    categories: state => state.region.categories || [],
    statuses: state => state.region.statuses || []
  }),
  methods: {
    ...mapActions([
      "listRegion",
      "listRegionCategory",
      "listRegionStatus",
      "updateRegion"
    ]),
    handleSearch() {
      this.query.offset = 0;
      this.fetch();
    },
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
    // async modify(data) {
    //   // 使用简单的方式，不修改路由参数
    //   this.mode = modifyMode;
    //   this.originRegion = Object.assign({}, data);
    //   this.currentRegion = data;
    // },
    // async update() {
    //   const { currentRegion, originRegion } = this;
    //   const update = diff(currentRegion, originRegion);
    //   if (update.modifiedCount === 0) {
    //     this.$message.warning("请修改要更新的信息");
    //     return;
    //   }
    //   try {
    //     await this.updateRegion({
    //       id: currentRegion.id,
    //       data: update.data
    //     });
    //     this.goBack();
    //   } catch (err) {
    //     this.$message.error(err.message);
    //   }
    // }
  },
  async beforeMount() {
    try {
      await this.listRegionCategory();
      await this.listRegionStatus();
    } catch (err) {
      this.$message.error(err.message);
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
.submit, .selector, .addBtn
  width: 100%
.addBtn
  margin-top: 15px
.pagination
  text-align: right
  margin-top: 15px
</style>
