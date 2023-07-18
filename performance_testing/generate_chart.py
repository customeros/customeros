import matplotlib.pyplot as plt
import pandas as pd
import numpy as np

# Load data from CSV file
data = pd.read_csv('average_response_times.csv')

# Extract the last 10 rows from the data
last_10_rows = data.iloc[-10:]

# Extract x-axis values from the first column
x = last_10_rows.iloc[:, 0]

# Extract y-axis values from all subsequent columns
y = last_10_rows.iloc[:, 1:]

# Set the figure size and margins
fig.subplots_adjust(bottom=0.2)  # Adjust the bottom margin

# Generate the chart
for i, column in enumerate(y.columns):
    x = np.arange(len(y[column]))
    plt.plot(x, y[column], marker='o', label=column)

# Set x-axis tick labels to the values from the first column
plt.xticks(np.arange(len(last_10_rows)), last_10_rows.iloc[:, 0], rotation='vertical')

# Add legend to show the lines for each column
plt.legend(bbox_to_anchor=(1.05, 1), loc='upper left', borderaxespad=0)

# Add measurement units to the x-axis and y-axis labels
plt.xlabel('Response time (ms)')
plt.ylabel('Request')

# Save the chart as an image
plt.savefig('average_response_times.png')