<template>
  <div class="region">
    <el-cascader
      clearable
      :props="props"
      @change="handleChange"
      v-model="region"
      :show-all-levels="$props.showAllLevels"
    ></el-cascader>
  </div>
</template>
<script>
import { mapActions } from "vuex";

export default {
  name: "Region",
  props: {
    maxLevel: {
      type: Number,
      default: 4
    },
    startLevel: {
      type: Number,
      default: 0
    },
    showAllLevels: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      region: "",
      query: {
        limit: 500,
        fields: "code,name"
      },
      props: {
        lazy: true,
        lazyLoad: this.lazyLoad
      }
    };
  },
  methods: {
    ...mapActions([
      // "listRegionCountry",
      // "listRegionProvince",
      // "listRegionCity",
      // "listRegionArea",
      // "listRegionStreet"
    ]),
    handleChange(value) {
      console.dir(value);
    },
    async lazyLoad(node, resolve) {
      console.dir(node);
      console.dir(resolve);
      //   const { maxLevel, startLevel } = this.$props;
      //   const { query } = this;
      //   let level = node.level + startLevel;
      //   let data = null;
      //   const leaf = level === maxLevel - 1;
      //   try {
      //     switch (level) {
      //       case countryLevel: {
      //         const { countries } = await this.listRegionCountry({
      //           params: query
      //         });
      //         data = countries;
      //         break;
      //       }
      //       case provinceLevel: {
      //         const { provinces } = await this.listRegionProvince({
      //           country: node.value || chinaCode,
      //           params: query
      //         });
      //         data = provinces;
      //         break;
      //       }
      //       case cityLevel: {
      //         const { cities } = await this.listRegionCity({
      //           province: node.value,
      //           params: query
      //         });
      //         data = cities;
      //         break;
      //       }
      //       case areaLevel: {
      //         const { areas } = await this.listRegionArea({
      //           city: node.value,
      //           params: query
      //         });
      //         data = areas;
      //         break;
      //       }
      //       case streetLevel: {
      //         const { streets } = await this.listRegionStreet({
      //           area: node.value,
      //           params: query
      //         });
      //         data = streets;
      //         break;
      //       }
      //       default:
      //         break;
      //     }
      //     if (!data) {
      //       throw new Error("获取数据失败");
      //     }
      //     const nodes = data.map(item => {
      //       let isLeaf = leaf;
      //       if (level === countryLevel && item.code !== chinaCode) {
      //         isLeaf = true;
      //       }
      //       return {
      //         value: item.code,
      //         label: item.name,
      //         leaf: isLeaf
      //       };
      //     });
      //     resolve(nodes);
      //   } catch (err) {
      //     this.$message.error(err.message);
      //   }
    }
  }
};
</script>
<style lang="sass" scoped>
.region
  .el-cascader
    width: 100%
</style>
