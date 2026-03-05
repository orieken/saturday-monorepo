<template>
  <div class="max-w-4xl mx-auto">
    <h2 class="text-3xl font-magical text-dnd-parchment mb-8 border-b border-dnd-slate pb-4">Your Loot Stash</h2>
    
    <div v-if="items.length" class="space-y-6">
      <div v-for="e in items" :key="e.id" class="bg-dnd-slate border border-gray-700 rounded-lg p-6 flex flex-col md:flex-row items-center justify-between shadow-lg" data-testid="cart-item">
        <div class="flex-1 mb-4 md:mb-0">
          <h3 class="text-xl font-magical text-dnd-gold mb-1" data-testid="cart-item-name">{{ e.name }}</h3>
          <div class="text-gray-400 text-sm" data-testid="cart-item-price">{{ e.priceSnapshot }} gp each</div>
        </div>
        
        <div class="flex items-center gap-6">
          <div class="flex items-center bg-gray-900 rounded-lg border border-gray-700 p-1">
            <button 
              @click="cart.setQuantity(e.id, e.quantity - 1)"
              class="w-8 h-8 flex items-center justify-center text-gray-400 hover:text-white hover:bg-gray-800 rounded transition-colors"
              data-testid="decrease-qty-btn"
            >
              -
            </button>
            <span class="w-12 text-center font-bold text-dnd-parchment" data-testid="cart-item-qty">{{ e.quantity }}</span>
            <button 
              @click="cart.setQuantity(e.id, e.quantity + 1)"
              class="w-8 h-8 flex items-center justify-center text-gray-400 hover:text-white hover:bg-gray-800 rounded transition-colors"
              data-testid="increase-qty-btn"
            >
              +
            </button>
          </div>
          
          <div class="text-xl font-bold text-dnd-gold min-w-[100px] text-right" data-testid="cart-item-total">
            {{ (e.priceSnapshot * e.quantity).toFixed(2) }} gp
          </div>
          
          <button 
            @click="cart.removeItem(e.id)"
            class="text-gray-500 hover:text-dnd-crimson transition-colors p-2"
            title="Remove from cart"
            data-testid="remove-item-btn"
          >
            <span class="text-xl">×</span>
          </button>
        </div>
      </div>

      <div class="bg-dnd-dark border border-dnd-slate rounded-lg p-6 mt-8">
        <div class="flex justify-between items-center mb-6">
          <span class="text-xl text-gray-400">Total Value</span>
          <span class="text-3xl font-magical text-dnd-gold" data-testid="cart-total">{{ totalPrice.toFixed(2) }} gp</span>
        </div>
        
        <!-- Error Message -->
        <div v-if="checkoutError" class="mb-4 p-4 bg-red-900/30 border border-red-700 rounded text-red-300 text-sm">
          {{ checkoutError }}
        </div>
        
        <button 
          @click="checkout"
          :disabled="checkingOut"
          class="w-full bg-dnd-crimson hover:bg-red-700 text-white py-4 rounded font-magical text-xl transition-colors shadow-lg active:transform active:scale-[0.99] flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          data-testid="checkout-btn"
        >
          <span v-if="!checkingOut">Complete Transaction</span>
          <span v-else>Processing...</span>
          <span class="text-2xl">🪙</span>
        </button>
        
        <p v-if="!authStore.isAuthenticated" class="mt-4 text-sm text-gray-400 text-center">
          You'll be asked to login before completing your purchase
        </p>
      </div>
    </div>
    
    <div v-else class="text-center py-20 bg-dnd-slate/30 rounded-lg border border-dnd-slate border-dashed">
      <div class="text-6xl mb-4">🕸️</div>
      <h3 class="text-2xl font-magical text-gray-400 mb-2" data-testid="cart-empty">Your stash is empty</h3>
      <p class="text-gray-500 mb-8" data-testid="cart-empty-message">Go forth and acquire some magical artifacts!</p>
      <router-link
        to="/"
        class="inline-block bg-dnd-slate hover:bg-gray-700 text-dnd-parchment px-6 py-3 rounded border border-gray-600 transition-colors"
        data-testid="return-to-shop-btn"
      >
        Return to Shop
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useCartStore } from '../stores/cartStore';
import { useAuthStore } from '../stores/authStore';
import { createOrder } from '../api/orders';

const router = useRouter();
const cart = useCartStore();
const authStore = useAuthStore();
const { itemsArray, totalPrice } = storeToRefs(cart);
const items = itemsArray;

const checkingOut = ref(false);
const checkoutError = ref<string | null>(null);

async function checkout() {
  // Check if user is authenticated
  if (!authStore.isAuthenticated) {
    // Redirect to login page
    router.push('/login');
    return;
  }

  checkingOut.value = true;
  checkoutError.value = null;

  try {
    // Create order with user's address if available
    const order = await createOrder(
      items.value,
      totalPrice.value,
      authStore.user?.address || undefined
    );

    console.log('Order created:', order.id);
    
    // Clear cart
    cart.clearCart();
    
    // Redirect to account page to view order
    router.push('/account?tab=orders');
  } catch (err) {
    console.error('Checkout failed', err);
    checkoutError.value = err instanceof Error ? err.message : 'Checkout failed. Please try again.';
  } finally {
    checkingOut.value = false;
  }
}
</script>
