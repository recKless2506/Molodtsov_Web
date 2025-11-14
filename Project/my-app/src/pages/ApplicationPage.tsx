// src/pages/ApplicationPage.tsx
import React from "react";
import type { Request } from "../types";
import "./application.css";

type Props = {
  requests: Request[];
  clearCart: () => void;
};

const ApplicationPage: React.FC<Props> = ({ requests, clearCart }) => {
  if (!requests || requests.length === 0) {
    return <div style={{ padding: 40 }}>Ваша корзина пуста.</div>;
  }
  const first = requests[0];
  return (
    <div className="page-container">
      <button onClick={clearCart} className="clear-button">Очистить корзину</button>

      <div className="input-container">
        <div className="input-block">
          <div className="input-label">Площадь помещения</div>
          <input className="input-field" value={String(first.PlaceSquare ?? "")} readOnly />
        </div>
        <div className="input-block">
          <div className="input-label">Температура за помещением</div>
          <input className="input-field" value={String(first.OutsideTemperature ?? "")} readOnly />
        </div>
        <div className="input-block">
          <div className="input-label">Температура в помещении</div>
          <input className="input-field" value={String(first.InsideTemperature ?? "")} readOnly />
        </div>
      </div>

      <div className="cards-container">
        {requests.map((r) => r.RequestHeaters?.map((rh) => (
          <div key={`${r.ID}-${rh.HeaterProduct.ID}`} className="card">
            <img className="card-img" src={rh.HeaterProduct.Image ?? "/static/images/teplodar_sputnik_elektro_6_updated.jpg"} alt={rh.HeaterProduct.Title} />
            <div className="card-text-row">
              <div className="card-text-block block-title">{rh.HeaterProduct.Title}</div>
              <div className="card-divider" />
              <div className="card-text-block block-specs">{rh.HeaterProduct.Power}</div>
              <div className="card-divider" />
              <div className="card-text-block block-input">
                <div className="input-label">Объём носителя</div>
                <input className="card-input" value={String(rh.Area ?? "")} readOnly style={{ color: "#000" }} />
              </div>
            </div>
          </div>
        )))}
      </div>
    </div>
  );
};

export default ApplicationPage;
