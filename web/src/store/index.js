import Vue from "vue";
import Vuex from "vuex";

import userStore from "@/store/modules/user";
import configStore from "@/store/modules/config";
import commonStore from "@/store/modules/common";
import brandStore from "@/store/modules/brand";
import productStore from "@/store/modules/product";
import regionStore from "@/store/modules/region";
import productCategoryStore from "@/store/modules/productCategory";
import supplierStore from "@/store/modules/supplier";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    user: userStore,
    config: configStore,
    common: commonStore,
    brand: brandStore,
    product: productStore,
    region: regionStore,
    productCategory: productCategoryStore,
    supplier: supplierStore
  }
});
