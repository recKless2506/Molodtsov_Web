import React from "react";
import Header from "../components/Header";
import ProductCard from "../components/ProductCard";
import type { HeaterProduct } from "../types";
import "./catalog.css";

type Props = {
  products: HeaterProduct[];
  cartCount: number;
  onAddToCart: (p: HeaterProduct) => void;
};

const CatalogPage: React.FC<Props> = ({ products, cartCount, onAddToCart }) => {
  return (
    <div>
      <Header cartCount={cartCount} />

      <div className="products">
        {products.length > 0 ? (
          products.map((product) => (
            <ProductCard
              key={product.ID}
              product={product}
              onAddToCart={onAddToCart}
            />
          ))
        ) : (
          <p style={{ textAlign: "center", width: "100%", marginTop: 40 }}>
            Товары не найдены.
          </p>
        )}
      </div>
    </div>
  );
};

export default CatalogPage;
