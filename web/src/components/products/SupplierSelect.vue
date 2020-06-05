<template>
  <el-select
    class="select"
    @change="handleChange"
    filterable
    remote
    reserve-keyword
    v-model="supplier"
    placeholder="请选择供应商"
    :remote-method="fetch"
    :loading="processing"
  >
    <el-option
      v-for="item in suppliers"
      :key="item.id"
      :label="`${item.name}`"
      :value="item.id"
    />
  </el-select>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { STATUS_ENABLED } from "@/constants/common";

export default {
  name: "SupplierSelect",
  props: {
    value: Number
  },
  data() {
    return {
      supplier: this.$props.value || null
    };
  },
  computed: {
    ...mapState({
      suppliers: state => state.supplier.list.data || [],
      processing: state => state.supplier.processing
    })
  },
  methods: {
    ...mapActions(["listSupplier"]),
    async fetch(query) {
      await this.listSupplier({
        limit: 20,
        status: STATUS_ENABLED,
        keyword: query
      });
    },
    handleChange(value) {
      this.$emit("input", value);
    }
  },

  beforeMount() {
    this.fetch();
  }
};
</script>
