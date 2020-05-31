<template>
  <el-select
    class="selector"
    @change="handleChange"
    filterable
    remote
    reserve-keyword
    v-model="brand"
    placeholder="请选择产品品牌"
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
import { STATUS_ENABLED } from "@/constants/common";

export default {
  name: "BrandSelect",
  props: {
    value: Number
  },
  data() {
    return {
      brand: this.$props.value || null
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
      this.$emit("input", value);
    },
    async fetch(query) {
      await this.listBrand({
        limit: 20,
        status: STATUS_ENABLED,
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
