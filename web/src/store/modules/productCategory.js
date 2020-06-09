import request from "@/helpers/request";
import { PRODUCT_CATEGORIES, PRODUCT_CATEGORIES_ID } from "@/constants/url";
import { attachStatusDesc, attachUpdatedAtDesc } from "@/store/modules/common";
import { addNoCacheQueryParam, findByID, toUploadFiles } from "@/helpers/util";

const prefix = "productCategory";
const mutationProductCategoryList = `${prefix}.list`;
const mutationProductCategoryListProcessing = `${mutationProductCategoryList}.processing`;
const mutationProductCategoryUpdate = `${prefix}.update`;

const state = {
  processing: false,
  list: {
    data: null,
    count: -1
  }
};

function enhanceBrandInfo(item) {
  attachUpdatedAtDesc(item);
  attachStatusDesc(item);
  item.files = toUploadFiles(item.icon);
}

export default {
  state,
  mutations: {
    [mutationProductCategoryListProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationProductCategoryList](state, { productCategories = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      productCategories.forEach(enhanceBrandInfo);
      state.list.data = productCategories;
    },
    [mutationProductCategoryUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
        enhanceBrandInfo(found);
      }
    }
  },
  actions: {
    // addProductCategory 添加产品分类
    async addProductCategory({ commit }, productCategory) {
      if (productCategory.files) {
        productCategory.icon = productCategory.files[0].url;
        delete productCategory.files;
      }
      commit(mutationProductCategoryListProcessing, true);
      try {
        const { data } = await request.post(
          PRODUCT_CATEGORIES,
          productCategory
        );
        return data;
      } finally {
        commit(mutationProductCategoryListProcessing, false);
      }
    },
    async listProductCategory({ commit }, params) {
      commit(mutationProductCategoryListProcessing, true);
      try {
        const { data } = await request.get(PRODUCT_CATEGORIES, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationProductCategoryList, data);
        return data;
      } finally {
        commit(mutationProductCategoryListProcessing, false);
      }
    },
    // getProductCategoryByID 通过id获取产品分类信息
    async getProductCategoryByID({ commit }, id) {
      const found = findByID(state.list.data, id);
      if (found) {
        return found;
      }
      commit(mutationProductCategoryListProcessing, true);
      try {
        const url = PRODUCT_CATEGORIES_ID.replace(":id", id);
        const { data } = await request.get(url, {
          params: addNoCacheQueryParam()
        });
        enhanceBrandInfo(data);
        return data;
      } finally {
        commit(mutationProductCategoryListProcessing, false);
      }
    },
    // updateProductCategoryByID 通过ID更新产品信息
    async updateProductCategoryByID({ commit }, { id, data }) {
      if (data.files) {
        data.icon = data.files[0].url;
        delete data.files;
      }
      commit(mutationProductCategoryListProcessing, true);
      try {
        const url = PRODUCT_CATEGORIES_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationProductCategoryUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationProductCategoryListProcessing, false);
      }
    }
  }
};
