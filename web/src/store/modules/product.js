import request from "@/helpers/request";

import { PRODUCTS, PRODUCTS_ID } from "@/constants/url";
import { formatDate, addNoCacheQueryParam, findByID } from "@/helpers/util";
import {
  listStatus,
  attachStatusDesc,
  attachUpdatedAtDesc
} from "@/store/modules/common";

const prefix = "product";
const mutationProductProcessing = `${prefix}.processing`;
const mutationProductList = `${prefix}.list`;
const mutationProductUpdate = `${prefix}.update`;

const state = {
  processing: false,
  list: {
    data: null,
    count: -1
  }
};

function fillAndUpdate(item) {
  if (!item.categories) {
    item.categories = [];
  }
  item.startedAtDesc = formatDate(item.startedAtDesc);
  item.endedAtDesc = formatDate(item.endedAt);
  attachUpdatedAtDesc(item);
  attachStatusDesc(item);
}

export default {
  state,
  mutations: {
    [mutationProductProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationProductList](state, { products = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      products.forEach(fillAndUpdate);
      state.list.data = products;
    },
    [mutationProductUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
        fillAndUpdate(found);
      }
    }
  },
  actions: {
    // addProduct 添加产品
    async addProduct({ commit }, product) {
      commit(mutationProductProcessing, true);
      try {
        const { data } = await request.post(PRODUCTS, product);
        return data;
      } finally {
        commit(mutationProductProcessing, false);
      }
    },
    listProductStatus: listStatus,
    // listProduct 获取产品
    async listProduct({ commit }, params) {
      commit(mutationProductProcessing, true);
      try {
        await listStatus({ commit });
        const { data } = await request.get(PRODUCTS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationProductList, data);
      } finally {
        commit(mutationProductProcessing, false);
      }
    },
    // getProductByID get product by id
    async getProductByID({ commit }, id) {
      const found = findByID(state.list.data, id);
      if (found) {
        return found;
      }
      commit(mutationProductProcessing, true);
      try {
        const url = PRODUCTS_ID.replace(":id", id);
        const { data } = await request.get(url, {
          params: addNoCacheQueryParam()
        });
        return data;
      } finally {
        commit(mutationProductProcessing, false);
      }
    },
    // updateProductByID update product by id
    async updateProductByID({ commit }, { id, data }) {
      if (!data || Object.keys(data).length === 0) {
        return;
      }
      commit(mutationProductProcessing, true);
      try {
        const url = PRODUCTS_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationProductUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationProductProcessing, false);
      }
    }
  }
};
