<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新产品分类信息"
    icon="el-icon-set-up"
    :id="categoryID"
    :findByID="getProductCategoryByID"
    :updateByID="updateProductCategoryByID"
    :add="addProductCategory"
    :fields="fields"
  />
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";

const categoryStatuses = [];
const fields = [
  {
    label: "名称：",
    key: "name",
    clearable: true,
    span: 6,
    placeholder: "请输入分类名称"
  },
  {
    label: "级别：",
    key: "level",
    type: "select",
    span: 6,
    placeholder: "请选择分类级别",
    options: [
      {
        name: "第一级",
        value: 1
      },
      {
        name: "第二级",
        value: 2
      },
      {
        name: "第三级",
        value: 3
      }
    ]
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    span: 6,
    options: categoryStatuses
  },
  {
    label: "上级分类：",
    labelWidth: "90px",
    key: "belongs",
    span: 6,
    type: "productCategory"
  }
];

export default {
  name: "ProductCategory",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      categoryID: 0,
      processing: false
    };
  },
  methods: {
    ...mapActions([
      "listProductCategoryStatus",
      "getProductCategoryByID",
      "updateProductCategoryByID",
      "addProductCategory"
    ])
  },
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.categoryID = Number(id);
    }
    try {
      const { statuses } = await this.listProductCategoryStatus();
      categoryStatuses.length = 0;
      categoryStatuses.push(...statuses);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
