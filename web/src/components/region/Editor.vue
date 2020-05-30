<template>
  <div class="brand">
    <BaseEditor
      v-if="!processing && fields"
      title="更新地区信息"
      icon="el-icon-position"
      :id="regionID"
      :findByID="getRegionByID"
      :updateByID="updateRegionByID"
      :fields="fields"
    />
  </div>
</template>
<script>
import BaseEditor from "@/components/base/Editor.vue";
import { mapActions } from "vuex";
const brandStatuses = [];
const fields = [
  {
    label: "名称：",
    key: "name"
  },
  {
    label: "代码：",
    key: "code",
    disabled: true
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    options: brandStatuses
  }
];

export default {
  name: "Brand",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      processing: false,
      regionID: 0
    };
  },
  methods: mapActions(["listBrandStatus", "updateRegionByID", "getRegionByID"]),
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.regionID = Number(id);
    }
    try {
      const { statuses } = await this.listBrandStatus();
      brandStatuses.length = 0;
      brandStatuses.push(...statuses);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
