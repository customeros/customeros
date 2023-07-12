import matplotlib.pyplot as plt
import pandas as pd

# Load data from CSV file
data = pd.read_csv('average_response_times.csv')

# Extract x-axis values from the first column
x = data.iloc[:, 0]

# Extract y-axis values from all subsequent columns
y = data.iloc[:, 1:]

# Generate the chart
for column in y.columns:
    plt.plot(x, y[column])

# Save the chart as an image
plt.savefig('average_response_times.png')
