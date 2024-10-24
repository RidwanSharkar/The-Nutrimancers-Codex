declare module 'react-liquid-gauge' {
        import * as React from 'react';
      
        export interface LiquidGaugeProps {
                style;
          value: number;
          width?: number;
          height?: number;
          min?: number;
          max?: number;
          forceRender?: boolean;
          textRenderer?: () => React.ReactNode;
          valueFormatter?: (value: number) => string;
          riseAnimation?: boolean;
          waveAnimation?: boolean;
          waveFrequency?: number;
          waveAmplitude?: number;
          gradient?: boolean;
          gradientStops?: {
            key: string;
            stopColor: string;
            stopOpacity: number;
            offset: string;
          }[];
          circleStyle?: React.CSSProperties;
          waveStyle?: React.CSSProperties;
          textStyle?: {
            fill: string;
            fontFamily: string;
            fontSize: string;
          };
          waveTextStyle?: {
            fill: string;
            fontFamily: string;
            fontSize: string;
          };
          onClick?: () => void;
          waveRenderer?: (props: { path: string }) => React.ReactNode;
        }
      
        const LiquidGauge: React.FC<LiquidGaugeProps>;
      
        export default LiquidGauge;
      }
      