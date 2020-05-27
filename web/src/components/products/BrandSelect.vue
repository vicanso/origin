<template>
  <el-select
    class="selector"
    @change="handleChange"
    v-model="brand"
    filterable
    remote
    reserve-keyword
    placeholder="请输入关键词"
    :remote-method="fetch"
    :loading="processing"
  >
    <el-option
      v-for="item in brands"
      :key="item.name"
      :label="item.name"
      :value="item.id"
    >
    </el-option>
  </el-select>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { BRAND_ENABLE } from "@/constants/common";

export default {
  name: "BrandSelect",
  data() {
    return {
      brand: null
    };
  },
  computed: {
    ...mapState({
      brands: state => state.brand.list.data || [],
      processing: state => state.brand.processing
    })
  },
  methods: {
    ...mapActions(["listBrand"]),
    handleChange(value) {
      this.$emit("change", value);
    },
    async fetch(query) {
      await this.listBrand({
        limit: 20,
        status: BRAND_ENABLE,
        keyword: query
      });
    }
  },
  beforeMount() {
    this.fetch();
  }
};
</script>
<style lang="sass" scoped>
.selector
  width: 100%
</style>
