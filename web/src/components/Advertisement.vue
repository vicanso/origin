<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新广告信息"
    icon="el-icon-c-scale-to-original"
    :id="id"
    :findByID="getAdvertisementByID"
    :updateByID="updateAdvertisementByID"
    :fields="fields"
    :add="addAdvertisement"
  />
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";

const advertisementStatuses = [];
const advertisementCategories = [];
const fields = [
  {
    label: "分类：",
    key: "category",
    type: "select",
    placeholder: "请选择广告分类",
    options: advertisementCategories,
    rules: [
      {
        required: true,
        message: "广告分类不能为空"
      }
    ]
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    placeholder: "请选择广告状态",
    options: advertisementStatuses,
    rules: [
      {
        required: true,
        message: "广告状态不能为空"
      }
    ]
  },
  {
    label: "链接：",
    key: "link",
    placeholder: "请配置广告链接",
    rules: [
      {
        required: true,
        message: "广告链接不能为空"
      }
    ]
  },
  {
    label: "简介：",
    key: "summary",
    placeholder: "请配置广告简介",
    rules: [
      {
        required: true,
        message: "广告简介不能为空"
      }
    ]
  },
  {
    label: "开始时间：",
    key: "startedAt",
    type: "datePicker",
    pickerType: "datetime",
    placeholder: "请选择广告生效时间",
    labelWidth: "100px",
    rules: [
      {
        required: true,
        message: "广告生效时间不能为空"
      }
    ]
  },
  {
    label: "结束时间：",
    key: "endedAt",
    type: "datePicker",
    pickerType: "datetime",
    placeholder: "请选择广告失效时间",
    labelWidth: "100px",
    rules: [
      {
        required: true,
        message: "广告失效时间不能为空"
      }
    ]
  },
  {
    label: "图片：",
    key: "files",
    span: 24,
    type: "upload",
    rules: [
      {
        required: true,
        message: "广告图片不能为空"
      }
    ]
  }
];

export default {
  name: "Advertisement",
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
      "addAdvertisement",
      "getAdvertisementByID",
      "updateAdvertisementByID",
      "listAdvertisementsCategory"
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
      advertisementStatuses.length = 0;
      advertisementStatuses.push(...statuses);
      const { categories } = await this.listAdvertisementsCategory();
      console.dir(categories);
      advertisementCategories.length = 0;
      advertisementCategories.push(...categories);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
