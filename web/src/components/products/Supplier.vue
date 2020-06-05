<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新供应商信息"
    icon="el-icon-good"
    :id="id"
    :findByID="getSupplierByID"
    :updateByID="updateSupplierByID"
    :fields="fields"
    :add="addSupplier"
  />
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";

const supplierStatuses = [];
const fields = [
  {
    label: "名称：",
    key: "name",
    clearable: true,
    placeholder: "请输入供应商名称",
    rules: [
      {
        required: true,
        message: "供应商名称不能为空"
      }
    ]
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    placeholder: "请选择供应商状态",
    options: supplierStatuses,
    rules: [
      {
        required: true,
        message: "供应商状态不能为空"
      }
    ]
  },
  {
    label: "联系人：",
    key: "contact",
    placeholder: "请输入联系人",
    rules: [
      {
        required: true,
        message: "供应商联系人不能为空"
      }
    ]
  },
  {
    label: "联系方式：",
    key: "mobile",
    placeholder: "请输入联系人电话",
    labelWidth: "100px",
    rules: [
      {
        required: true,
        message: "供应商联系方式不能为空"
      }
    ]
  },
  {
    label: "地址：",
    key: "baseAddress",
    type: "region",
    placeholder: "请选择供应商地址",
    rules: [
      {
        required: true,
        message: "供应商地址不能为空"
      }
    ]
  },
  {
    label: "",
    key: "address",
    labelWidth: "0px",
    placeholder: "请输入供应商地址",
    rules: [
      {
        required: true,
        message: "供应商地址不能为空"
      }
    ]
  }
];

export default {
  name: "Supplier",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      id: 0,
      processing: false
    };
  },
  methods: {
    ...mapActions([
      "listStatus",
      "addSupplier",
      "getSupplierByID",
      "updateSupplierByID"
    ])
  },
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.id = Number(id);
    }
    try {
      const { statuses } = await this.listStatus();
      supplierStatuses.length = 0;
      supplierStatuses.push(...statuses);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
