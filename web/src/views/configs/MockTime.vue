<template>
  <ConfigEditor
    name="添加/更新MockTime配置"
    summary="针对应用时间Mock，用于测试环境中调整应用时间"
    :category="category"
    :defaultValue="defaultValue"
    v-if="!processing"
    :id="currentID"
  >
    <template v-slot:data="configProps">
      <MockTimeData
        :data="configProps.form.data"
        @change="configProps.form.data = $event"
      />
    </template>
  </ConfigEditor>
</template>
<script>
import ConfigEditor from "@/components/configs/Editor.vue";
import MockTimeData from "@/components/configs/MockTimeData.vue";
import { MOCK_TIME } from "@/constants/config";
import { mapActions } from "vuex";

export default {
  name: "MockTime",
  components: {
    MockTimeData,
    ConfigEditor
  },
  data() {
    return {
      defaultValue: {
        name: MOCK_TIME,
        category: MOCK_TIME
      },
      processing: false,
      currentID: 0,
      category: MOCK_TIME
    };
  },
  methods: {
    ...mapActions(["listConfig"])
  },
  async beforeMount() {
    this.processing = true;
    try {
      const { configs } = await this.listConfig({
        name: MOCK_TIME
      });
      if (configs.length !== 0) {
        this.$router.replace({
          query: {
            id: configs[0].id
          }
        });
      }
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
