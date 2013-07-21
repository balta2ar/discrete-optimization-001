from Numberjack import *

def model_warehouse_planning(data):
    WareHouseOpen = VarArray(data.NumberOfWarehouses)

    ShopSupplied = Matrix(data.NumberOfWarehouses,
                          data.NumberOfShops)

    # Cost of running warehouses
    warehouseCost = Sum(WareHouseOpen, data.WareHouseCosts)

    # Cost of shops using warehouses
    transpCost = Sum([ Sum(varRow, costRow)
                       for (varRow, costRow) in zip(ShopSupplied, data.SupplyCost)])

    obj = warehouseCost + transpCost

    model = Model(
        # Objective function
        Minimise(obj),
        # Channel from store opening to store supply matrix
        [[var <= store for var in col]
         for (col, store) in zip(ShopSupplied.col, WareHouseOpen)],
        # Make sure every shop if supplied by one store
        [Sum(row) == 1 for row in ShopSupplied.row],
        # Make sure that each store does not exceed it's supply capacity
        [Sum(col) <= cap
         for (col, cap) in zip(ShopSupplied.col, data.Capacity)]
    )

    return (obj, WareHouseOpen, ShopSupplied, model)

def solve_warehouse_planning(data, param):
    (obj, WareHouseOpen, ShopSupplied, model) = model_warehouse_planning(data)
    solver = model.load(param['solver'])
    solver.setVerbosity(1)
    solver.solve()
    print obj.get_value()
    print "",WareHouseOpen
    print ShopSupplied

class WareHouseData:
    def __init__(self):
        self.NumberOfWarehouses = 5
        self.NumberOfShops = 10
        self.FixedCost = 30
        self.WareHouseCosts = [30, 30, 30, 30, 30]
        self.Capacity = [1,4,2,1,3]
        self.SupplyCost = supplyCost = [
            [ 20, 24, 11, 25, 30 ],
            [ 28, 27, 82, 83, 74 ],
            [ 74, 97, 71, 96, 70 ],
            [ 2, 55, 73, 69, 61 ],
            [ 46, 96, 59, 83, 4 ],
            [ 42, 22, 29, 67, 59 ],
            [ 1, 5, 73, 59, 56 ],
            [ 10, 73, 13, 43, 96 ],
            [ 93, 35, 63, 85, 46 ],
            [ 47, 65, 55, 71, 95 ]
        ]

solve_warehouse_planning(WareHouseData(), input({'solver':'SCIP'}))
