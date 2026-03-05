<template>
  <div class="max-w-6xl mx-auto">
    <!-- Page Header -->
    <div class="mb-8 border-b border-dnd-slate pb-4">
      <h2 class="text-3xl font-magical text-dnd-parchment">Your Magical Account</h2>
      <p class="text-gray-400 mt-2">Manage your profile and view your order history</p>
    </div>

    <!-- Tabs -->
    <div class="mb-8 flex gap-4 border-b border-dnd-slate">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="[
          'px-6 py-3 font-magical text-lg transition-all relative',
          activeTab === tab.id
            ? 'text-dnd-gold border-b-2 border-dnd-gold'
            : 'text-gray-400 hover:text-dnd-parchment'
        ]"
        :data-testid="tab.id === 'orders' ? 'order-history-tab' : `tab-${tab.id}`"
      >
        {{ tab.icon }} {{ tab.label }}
      </button>
    </div>

    <!-- Profile Tab -->
    <div v-if="activeTab === 'profile'" class="space-y-6">
      <div class="bg-dnd-dark border-2 border-dnd-slate rounded-lg p-6">
        <h3 class="text-xl font-magical text-dnd-gold mb-4">Profile Information</h3>
        
        <div v-if="!editingProfile" class="space-y-4">
          <div>
            <label class="text-sm text-gray-400">Name</label>
            <p class="text-lg text-dnd-parchment">{{ authStore.user?.name }}</p>
          </div>
          <div>
            <label class="text-sm text-gray-400">Email</label>
            <p class="text-lg text-dnd-parchment">{{ authStore.user?.email }}</p>
          </div>
          <div>
            <label class="text-sm text-gray-400">Member Since</label>
            <p class="text-lg text-dnd-parchment">{{ formatDate(authStore.user?.createdAt) }}</p>
          </div>
          
          <button
            @click="editingProfile = true"
            class="mt-4 px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors"
            data-testid="edit-profile-btn"
          >
            ✏️ Edit Profile
          </button>
        </div>

        <form v-else @submit.prevent="handleUpdateProfile" class="space-y-4">
          <div>
            <label for="profile-name" class="block text-sm text-gray-400 mb-2">Name</label>
            <input
              id="profile-name"
              v-model="profileForm.name"
              type="text"
              required
              class="w-full px-4 py-2 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold outline-none text-white"
              data-testid="profile-name-input"
            />
          </div>

          <div class="flex gap-4">
            <button
              type="submit"
              :disabled="authStore.loading"
              class="px-6 py-2 bg-green-600 hover:bg-green-700 text-white rounded transition-colors disabled:opacity-50"
              data-testid="save-profile-btn"
            >
              💾 Save Changes
            </button>
            <button
              type="button"
              @click="cancelEdit"
              class="px-6 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded transition-colors"
              data-testid="cancel-edit-btn"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>

      <!-- Shipping Address -->
      <div class="bg-dnd-dark border-2 border-dnd-slate rounded-lg p-6">
        <h3 class="text-xl font-magical text-dnd-gold mb-4">Shipping Address</h3>
        
        <div v-if="!editingAddress" class="space-y-2">
          <div v-if="authStore.user?.address">
            <p class="text-dnd-parchment">{{ authStore.user.address.street }}</p>
            <p class="text-dnd-parchment">
              {{ authStore.user.address.city }}, {{ authStore.user.address.state }} {{ authStore.user.address.zip }}
            </p>
            <p class="text-dnd-parchment">{{ authStore.user.address.country }}</p>
          </div>
          <p v-else class="text-gray-400 italic">No shipping address set</p>
          
          <button
            @click="startEditAddress"
            class="mt-4 px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded transition-colors"
            data-testid="edit-address-btn"
          >
            ✏️ {{ authStore.user?.address ? 'Edit' : 'Add' }} Address
          </button>
        </div>

        <form v-else @submit.prevent="handleUpdateAddress" class="space-y-4">
          <div>
            <label for="address-street" class="block text-sm text-gray-400 mb-2">Street Address</label>
            <input
              id="address-street"
              v-model="addressForm.street"
              type="text"
              required
              class="w-full px-4 py-2 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold outline-none text-white"
              data-testid="address-street-input"
            />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="address-city" class="block text-sm text-gray-400 mb-2">City</label>
              <input
                id="address-city"
                v-model="addressForm.city"
                type="text"
                required
                class="w-full px-4 py-2 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold outline-none text-white"
                data-testid="address-city-input"
              />
            </div>
            <div>
              <label for="address-state" class="block text-sm text-gray-400 mb-2">State/Province</label>
              <input
                id="address-state"
                v-model="addressForm.state"
                type="text"
                required
                class="w-full px-4 py-2 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold outline-none text-white"
                data-testid="address-state-input"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="address-zip" class="block text-sm text-gray-400 mb-2">ZIP Code</label>
              <input
                id="address-zip"
                v-model="addressForm.zip"
                type="text"
                required
                class="w-full px-4 py-2 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold outline-none text-white"
                data-testid="address-zip-input"
              />
            </div>
            <div>
              <label for="address-country" class="block text-sm text-gray-400 mb-2">Country</label>
              <input
                id="address-country"
                v-model="addressForm.country"
                type="text"
                required
                class="w-full px-4 py-2 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold outline-none text-white"
                data-testid="address-country-input"
              />
            </div>
          </div>

          <div class="flex gap-4">
            <button
              type="submit"
              :disabled="authStore.loading"
              class="px-6 py-2 bg-green-600 hover:bg-green-700 text-white rounded transition-colors disabled:opacity-50"
              data-testid="save-address-btn"
            >
              💾 Save Address
            </button>
            <button
              type="button"
              @click="cancelEditAddress"
              class="px-6 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded transition-colors"
              data-testid="cancel-address-btn"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Orders Tab -->
    <div v-if="activeTab === 'orders'" class="space-y-6">
      <div v-if="loadingOrders" class="text-center py-12">
        <p class="text-gray-400 text-lg">✨ Loading your order history...</p>
      </div>

      <div v-else-if="orders.length === 0" class="text-center py-12 bg-dnd-dark border-2 border-dnd-slate rounded-lg">
        <p class="text-gray-400 text-lg mb-4">📜 No orders yet</p>
        <router-link to="/" class="text-dnd-gold hover:text-yellow-400 transition-colors">
          Start shopping for magical items →
        </router-link>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="order in orders"
          :key="order.id"
          class="bg-dnd-dark border-2 border-dnd-slate rounded-lg p-6 hover:border-dnd-gold transition-colors"
          :data-testid="'order-item'"
        >
          <div class="flex justify-between items-start mb-4">
            <div>
              <h3 class="text-xl font-magical text-dnd-gold">Order #{{ order.id }}</h3>
              <p class="text-sm text-gray-400" data-testid="order-date">{{ formatDate(order.createdAt) }}</p>
            </div>
            <span
              :class="[
                'px-4 py-1 rounded text-sm font-medium',
                getStatusClass(order.status)
              ]"
            >
              {{ order.status.toUpperCase() }}
            </span>
          </div>

          <div class="space-y-2 mb-4">
            <div
              v-for="item in order.items"
              :key="item.id"
              class="flex justify-between text-dnd-parchment"
              data-testid="order-line-item"
            >
              <span>{{ item.name }} × {{ item.quantity }}</span>
              <span>{{ formatPrice(item.priceSnapshot * item.quantity) }}</span>
            </div>
          </div>

          <div class="border-t border-dnd-slate pt-4 flex justify-between items-center">
            <span class="text-lg font-magical text-dnd-gold">Total</span>
            <span class="text-xl font-bold text-dnd-gold" data-testid="order-total">{{ formatPrice(order.total) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Logout Button -->
    <div class="mt-8 pt-8 border-t border-dnd-slate">
      <button
        @click="handleLogout"
        class="px-6 py-2 bg-red-600 hover:bg-red-700 text-white rounded transition-colors"
        data-testid="logout-button"
      >
        🚪 Logout
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/authStore';
import { fetchOrders } from '../api/orders';
import type { Order } from '../api/orders';
import type { Address } from '../stores/authStore';

const router = useRouter();
const authStore = useAuthStore();

const activeTab = ref<'profile' | 'orders'>('profile');
const tabs = [
  { id: 'profile' as const, label: 'Profile', icon: '👤' },
  { id: 'orders' as const, label: 'Order History', icon: '📦' }
];

const editingProfile = ref(false);
const editingAddress = ref(false);
const loadingOrders = ref(false);
const orders = ref<Order[]>([]);

const profileForm = ref({
  name: ''
});

const addressForm = ref<Address>({
  street: '',
  city: '',
  state: '',
  zip: '',
  country: ''
});

onMounted(async () => {
  // Redirect if not authenticated
  if (!authStore.isAuthenticated) {
    router.push('/login');
    return;
  }

  // Check for tab query parameter
  const tabParam = router.currentRoute.value.query.tab as string;
  if (tabParam === 'orders') {
    activeTab.value = 'orders';
  }

  // Load orders
  await loadOrders();
});

async function loadOrders() {
  loadingOrders.value = true;
  try {
    orders.value = await fetchOrders();
    // Sort by date, newest first
    orders.value.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
  } catch (err) {
    console.error('Failed to load orders', err);
  } finally {
    loadingOrders.value = false;
  }
}

function startEditAddress() {
  if (authStore.user?.address) {
    addressForm.value = { ...authStore.user.address };
  } else {
    addressForm.value = {
      street: '',
      city: '',
      state: '',
      zip: '',
      country: ''
    };
  }
  editingAddress.value = true;
}

function cancelEdit() {
  editingProfile.value = false;
  profileForm.value.name = authStore.user?.name || '';
}

function cancelEditAddress() {
  editingAddress.value = false;
}

async function handleUpdateProfile() {
  const success = await authStore.updateProfile({ name: profileForm.value.name });
  if (success) {
    editingProfile.value = false;
  }
}

async function handleUpdateAddress() {
  const success = await authStore.updateProfile({ address: addressForm.value });
  if (success) {
    editingAddress.value = false;
  }
}

async function handleLogout() {
  await authStore.logout();
  router.push('/');
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return 'N/A';
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  });
}

function formatPrice(price: number): string {
  return `${price.toFixed(2)} gp`;
}

function getStatusClass(status: string): string {
  const classes = {
    pending: 'bg-yellow-900/30 text-yellow-300 border border-yellow-700',
    processing: 'bg-blue-900/30 text-blue-300 border border-blue-700',
    completed: 'bg-green-900/30 text-green-300 border border-green-700',
    cancelled: 'bg-red-900/30 text-red-300 border border-red-700'
  };
  return classes[status as keyof typeof classes] || classes.pending;
}

// Initialize profile form
if (authStore.user) {
  profileForm.value.name = authStore.user.name;
}
</script>
