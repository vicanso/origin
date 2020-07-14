import request from "@/helpers/request";
import { ORDERS, ORDERS_ID, ORDERS_STATUSES } from "@/constants/url";
import { findByID, formatDate } from "@/helpers/util";

const prefix = "order";
const mutationOrderProcessing = `${prefix}.processing`;
const mutationOrderList = `${prefix}.list`;
const mutationOrderUpdate = `${prefix}.update`;
const mutationOrderStatuses = `${prefix}.statuses`;

const state = {
  processing: false,
  list: {
    data: null,
    count: -1
  },
  statuses: null
};

function enhanceOrder(order) {
  order.address = order.receiverBaseAddressDesc + order.receiverAddress;
  order.payAmount = order.payAmount.toFixed(2);
  order.amount = order.amount.toFixed(2);
  order.createdAt = formatDate(order.createdAt);
}

export default {
  state,
  mutations: {
    [mutationOrderProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationOrderList](state, { orders = [], count = 0 }) {
      if (count >= 0) {
        state.list.count = count;
      }
      orders.forEach(enhanceOrder);
      state.list.data = orders;
    },
    [mutationOrderUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
        enhanceOrder(found);
      }
    },
    [mutationOrderStatuses](state, { statuses = [] }) {
      state.statuses = statuses;
    }
  },
  actions: {
    // listOrder get the order list
    async listOrder({ commit }, params) {
      commit(mutationOrderProcessing, true);
      try {
        const { data } = await request.get(ORDERS, {
          params
        });
        commit(mutationOrderList, data);
      } finally {
        commit(mutationOrderProcessing, false);
      }
    },
    // getOrderByID get order by id
    async getOrderBySN({ commit }, sn) {
      commit(mutationOrderProcessing, true);
      try {
        const url = ORDERS_ID.replace(":sn", sn);
        const { data } = await request.get(url);
        return data;
      } finally {
        commit(mutationOrderProcessing, false);
      }
    },
    async listOrderStatus({ commit }) {
      const { data } = await request.get(ORDERS_STATUSES);
      commit(mutationOrderStatuses, data);
      return data;
    }
  }
};
