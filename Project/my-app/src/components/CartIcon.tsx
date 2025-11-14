// src/components/CartIcon.tsx
import React from "react";
import "./carticon.css";

type CartIconProps = {
  count: number;
};

const CartIcon: React.FC<CartIconProps> = ({ count }) => {
  return (
    <div className="cart-icon-wrapper" title="Корзина">
      <div className="cart-icon">
        <svg viewBox="0 0 24 24">
          <path fill="white" d="M7 4h-2l-3 9v2h2l3-9zm0 11a2 2 0 1 0 4 0 2 2 0 0 0-4 0zm9-11h-6l-1 2h8v2h-8l1 2h6v2h-6l1 2h6v2h-6v-2h-6v-2h6l-1-2h-6v-2h6l-1-2h-6v-2h6l-1-2h6v-2z"/>
        </svg>
        {count > 0 && <div className="cart-count">{count}</div>}
      </div>
    </div>
  );
};

export default CartIcon;
