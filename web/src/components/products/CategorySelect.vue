<template>
  <el-select
    class="select"
    @change="handleChange"
    filterable
    remote
    reserve-keyword
    v-model="categories"
    multiple
    placeholder="请输入关键词"
    :remote-method="fetch"
    :loading="processing"
  >
    <el-option
      v-for="item in productCategories"
      :key="item.id"
      :label="`${item.name}(${item.level})`"
      :value="item.id"
    />
  </el-select>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { STATUS_ENABLED } from "@/constants/common";
export default {
  name: "ProductCategorySelect",
  props: {
    value: Array,
    level: Number
  },
  data() {
    console.dir(this.$props.level);
    return {
      categories: this.$props.value || []
    };
  },
  computed: {
    ...mapState({
      productCategories: state => state.productCategory.list.data || [],
      processing: state => state.productCategory.processing
    })
  },
  methods: {
    ...mapActions(["listProductCategory"]),
    handleChange(value) {
      this.$emit("input", value);
    },
    async fetch(query) {
      await this.listProductCategory({
        limit: 100,
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
