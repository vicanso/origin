import request from "@/helpers/request";

import { PRODUCTS_LIST_STATUS, PRODUCTS } from "@/constants/url";
import { formatDate, addNoCacheQueryParam } from "@/helpers/util";

const prefix = "product";
const mutationProductListStatus = `${prefix}.list.status`;
const mutationProductListStatusProcessing = `${mutationProductListStatus}.processing`;
const mutationProductProcessing = `${prefix}.processing`;
const mutationProductList = `${prefix}.list`;
// const mutationProductUpdate = `${prefix}.update`;

const state = {
  statusesListProcessing: false,
  statuses: null,
  processing: false,
  list: {
    data: null,
    count: -1
  }
};

async function listProductStatus({ commit }) {
  if (state.statuses) {
    return;
  }
  commit(mutationProductListStatusProcessing, true);
  try {
    const { data } = await request.get(PRODUCTS_LIST_STATUS);
    commit(mutationProductListStatus, data);
  } finally {
    commit(mutationProductListStatusProcessing, false);
  }
}

function fillAndUpdate(item) {
  if (!item.categories) {
    item.categories = [];
  }
  item.updatedAtDesc = formatDate(item.updatedAt);
  item.startedAtDesc = formatDate(item.startedAtDesc);
  item.endedAtDesc = formatDate(item.endedAt);
  state.statuses.forEach(status => {
    if (status.value === item.status) {
      item.statusDesc = status.name;
    }
  });
}

export default {
  state,
  mutations: {
    [mutationProductListStatusProcessing](state, processing) {
      state.statusesListProcessing = processing;
    },
    [mutationProductListStatus](state, { statuses }) {
      state.statuses = statuses;
    },
    [mutationProductProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationProductList](state, { products = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      products.forEach(fillAndUpdate);
      state.list.data = products;
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
    listProductStatus,
    // listProduct 获取产品
    async listProduct({ commit }, params) {
      commit(mutationProductProcessing, true);
      try {
        await listProductStatus({ commit });
        const { data } = await request.get(PRODUCTS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationProductList, data);
      } finally {
        commit(mutationProductProcessing, false);
      }
    }
  }
};
