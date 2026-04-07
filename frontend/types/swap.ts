export type SwapQuote = {
  chain: string;
  input_mint: string;
  output_mint: string;
  input_amount: string;
  output_amount: string;
  price_impact?: string;
  route?: string;
};
