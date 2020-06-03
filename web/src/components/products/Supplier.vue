<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新供应商信息"
    icon="el-icon-good"
    :id="supplierID"
    :findByID="getSupplierByID"
    :updateByID="updateSupplierByID"
    :fields="fields"
    :add="addSupplier"
    :rules="rules"
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
    placeholder: "请输入供应商名称"
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    placeholder: "请选择供应商状态",
    options: supplierStatuses
  },
  {
    label: "联系人：",
    key: "contact",
    placeholder: "请输入联系人"
  },
  {
    label: "联系方式：",
    key: "mobile",
    placeholder: "请输入联系人电话",
    labelWidth: "100px"
  },
  {
    label: "地址：",
    key: "address",
    placeholder: "请输入供应商地址",
    span: 16
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
      brandID: 0,
      processing: false,
      rules: {
        name: [
          {
            required: true,
            message: "供应商名称不能为空"
          }
        ],
        status: [
          {
            required: true,
            message: "供应商状态不能为空"
          }
        ],
        contact: [
          {
            required: true,
            message: "供应商联系人不能为空"
          }
        ],
        mobile: [
          {
            required: true,
            message: "供应商联系方式不能为空"
          }
        ],
        address: [
          {
            required: true,
            message: "供应商地址不能为空"
          }
        ]
      }
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
      this.supplierID = Number(id);
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
