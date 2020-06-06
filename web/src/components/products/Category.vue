<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新产品分类信息"
    icon="el-icon-set-up"
    :id="id"
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
    placeholder: "请输入分类名称",
    rules: [
      {
        required: true,
        message: "产品分类名称不能为空"
      }
    ]
  },
  {
    label: "级别：",
    key: "level",
    type: "select",
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
    ],
    rules: [
      {
        required: true,
        message: "产品分类级别不能为空"
      }
    ]
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    options: categoryStatuses,
    rules: [
      {
        required: true,
        message: "产品分类状态不能为空"
      }
    ]
  },
  {
    label: "热度：",
    key: "hot",
    dataType: "number",
    placeholder: "请输入产品热度(1-1000)"
  },
  {
    label: "上级分类：",
    labelWidth: "90px",
    key: "belongs",
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
      id: 0,
      processing: false
    };
  },
  methods: {
    ...mapActions([
      "listStatus",
      "getProductCategoryByID",
      "updateProductCategoryByID",
      "addProductCategory"
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
