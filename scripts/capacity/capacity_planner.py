import requests
import json
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
import os

class CapacityPlanner:
    def __init__(self):
        self.prometheus_url = os.getenv('PROMETHEUS_URL', 'http://prometheus:9090')
        self.prediction_days = int(os.getenv('PREDICTION_DAYS', '30'))
        self.alert_threshold = float(os.getenv('ALERT_THRESHOLD', '0.8'))

    def query_prometheus(self, query, start_time, end_time, step='1h'):
        """Query Prometheus for range data"""
        params = {
            'query': query,
            'start': start_time.isoformat() + 'Z',
            'end': end_time.isoformat() + 'Z',
            'step': step
        }

        response = requests.get(f'{self.prometheus_url}/api/v1/query_range', params=params)
        return response.json()

    def predict_resource_usage(self, metric_name, query):
        """Predict future resource usage using linear regression"""
        end_time = datetime.utcnow()
        start_time = end_time - timedelta(days=7)  # Use last 7 days for prediction

        data = self.query_prometheus(query, start_time, end_time)

        if not data['data']['result']:
            return None

        # Extract time series data
        values = data['data']['result'][0]['values']
        timestamps = [float(v[0]) for v in values]
        metrics = [float(v[1]) for v in values if v[1] != 'NaN']

        if len(metrics) < 2:
            return None

        # Simple linear regression
        x = np.array(range(len(metrics)))
        y = np.array(metrics)

        # Calculate trend
        slope = np.polyfit(x, y, 1)[0]
        current_value = metrics[-1]

        # Predict future values
        future_hours = self.prediction_days * 24
        predicted_value = current_value + (slope * future_hours)

        return {
            'metric': metric_name,
            'current_value': current_value,
            'predicted_value': predicted_value,
            'trend_slope': slope,
            'prediction_days': self.prediction_days,
            'will_exceed_threshold': predicted_value > self.alert_threshold
        }

    def generate_capacity_report(self):
        """Generate comprehensive capacity planning report"""
        metrics = [
            {
                'name': 'CPU Usage',
                'query': '100 - (avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)',
                'unit': '%',
                'threshold': 80
            },
            {
                'name': 'Memory Usage',
                'query': '(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100',
                'unit': '%',
                'threshold': 85
            },
            {
                'name': 'Disk Usage',
                'query': '(node_filesystem_size_bytes{mountpoint="/"} - node_filesystem_avail_bytes{mountpoint="/"}) / node_filesystem_size_bytes{mountpoint="/"} * 100',
                'unit': '%',
                'threshold': 90
            },
            {
                'name': 'API Request Rate',
                'query': 'rate(http_requests_total[5m])',
                'unit': 'req/s',
                'threshold': 1000
            }
        ]

        report = {
            'generated_at': datetime.utcnow().isoformat(),
            'prediction_period_days': self.prediction_days,
            'predictions': [],
            'recommendations': []
        }

        for metric in metrics:
            prediction = self.predict_resource_usage(metric['name'], metric['query'])

            if prediction:
                prediction['unit'] = metric['unit']
                prediction['threshold'] = metric['threshold']
                report['predictions'].append(prediction)

                # Generate recommendations
                if prediction['will_exceed_threshold']:
                    if 'CPU' in metric['name']:
                        report['recommendations'].append({
                            'metric': metric['name'],
                            'action': 'Scale horizontally or optimize CPU-intensive processes',
                            'priority': 'high' if prediction['predicted_value'] > 90 else 'medium'
                        })
                    elif 'Memory' in metric['name']:
                        report['recommendations'].append({
                            'metric': metric['name'],
                            'action': 'Increase memory allocation or optimize memory usage',
                            'priority': 'high' if prediction['predicted_value'] > 95 else 'medium'
                        })
                    elif 'Disk' in metric['name']:
                        report['recommendations'].append({
                            'metric': metric['name'],
                            'action': 'Expand storage or implement cleanup policies',
                            'priority': 'critical' if prediction['predicted_value'] > 95 else 'high'
                        })

        return report

    def save_report(self, report):
        """Save capacity planning report"""
        os.makedirs('/data', exist_ok=True)

        timestamp = datetime.utcnow().strftime('%Y%m%d_%H%M%S')
        filename = f'/data/capacity_report_{timestamp}.json'

        with open(filename, 'w') as f:
            json.dump(report, f, indent=2)

        # Also save as latest
        with open('/data/capacity_report_latest.json', 'w') as f:
            json.dump(report, f, indent=2)

        return filename

def main():
    planner = CapacityPlanner()

    print("🔮 Generating capacity planning report...")
    report = planner.generate_capacity_report()

    filename = planner.save_report(report)
    print(f"📊 Report saved: {filename}")

    # Print summary
    print("\n📋 Capacity Planning Summary:")
    print("=" * 40)

    for prediction in report['predictions']:
        status = "⚠️ ALERT" if prediction['will_exceed_threshold'] else "✅ OK"
        print(f"{prediction['metric']}: {prediction['current_value']:.1f}{prediction.get('unit', '')} → {prediction['predicted_value']:.1f}{prediction.get('unit', '')} {status}")

    if report['recommendations']:
        print(f"\n🎯 Recommendations ({len(report['recommendations'])}):")
        for rec in report['recommendations']:
            priority = rec['priority'].upper()
            print(f"[{priority}] {rec['metric']}: {rec['action']}")

if __name__ == "__main__":
    main()

