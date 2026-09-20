import { createRouter, createWebHistory } from 'vue-router';
import ShopHome from '../views/ShopHome.vue';
import ItemPage from '../views/ItemPage.vue';
import CartView from '../views/CartView.vue';
import LoginView from '../views/LoginView.vue';
import AccountView from '../views/AccountView.vue';
import LocationsView from '../views/LocationsView.vue';

const routes = [
  { path: '/', component: ShopHome },
  { path: '/item/:id', component: ItemPage, props: true },
  { path: '/cart', component: CartView },
  { path: '/login', component: LoginView },
  { path: '/account', component: AccountView },
  { path: '/locations', component: LocationsView }
];

export default createRouter({ history: createWebHistory(), routes });


