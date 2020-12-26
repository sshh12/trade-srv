from sklearn.ensemble import (
    RandomForestClassifier,
    RandomForestRegressor,
    AdaBoostClassifier,
    AdaBoostRegressor,
    GradientBoostingClassifier,
    GradientBoostingRegressor,
    ExtraTreesClassifier,
    ExtraTreesRegressor,
    VotingClassifier,
    VotingRegressor,
    StackingClassifier,
    StackingRegressor,
)
from sklearn.metrics import accuracy_score
from sklearn.tree import DecisionTreeRegressor
from sklearn.neighbors import KNeighborsRegressor
from sklearn.linear_model import BayesianRidge, LinearRegression, RidgeCV
from sklearn.model_selection import cross_val_score, train_test_split
from sklearn.neural_network import MLPRegressor
from xgboost import XGBRegressor
import re


class ModelRepo:
    def __init__(self, seed=None, n_jobs=None):
        self.seed = seed
        self.n_jobs = n_jobs

    def get_model(self, name):
        raise NotImplementedError()

    def iter_fit(self, X, y=None, models=[]):
        for m in models:
            if isinstance(m, str):
                model = self.get_model(m)
                name = m
            else:
                model = m
                name = str(m)
            model.fit(X, y)
            yield name, model


class ClassifierRepo(ModelRepo):
    def get_model(self, model):
        if match := re.match(r"ridgecv", model):
            return RidgeCV()
        elif match := re.match(r"rf(\d+)", model):
            return RandomForestClassifier(n_jobs=self.n_jobs, n_estimators=int(match.group(1)), random_state=self.seed)
        elif match := re.match(r"rf(\d+)d(\d+)", model):
            return RandomForestClassifier(
                n_jobs=self.n_jobs,
                n_estimators=int(match.group(1)),
                random_state=self.seed,
                max_depth=int(match.group(2)),
            )
        elif match := re.match(r"ada(\d+)", model):
            return AdaBoostClassifier(n_estimators=int(match.group(1)), random_state=self.seed)
        elif match := re.match(r"xtrees(\d+)", model):
            return ExtraTreesClassifier(n_estimators=int(match.group(1)), random_state=self.seed, n_jobs=self.n_jobs)
        elif match := re.match(r"vote:([\w,]+)", model):
            models = [(mn, self.get_model(mn)) for mn in match.group(1).split(",")]
            return VotingClassifier(models, n_jobs=self.n_jobs, voting="hard")
        elif match := re.match(r"softvote:([\w,]+)", model):
            models = [(mn, self.get_model(mn)) for mn in match.group(1).split(",")]
            return VotingClassifier(models, n_jobs=self.n_jobs, voting="soft")
        elif match := re.match(r"stack:(\w+):([\w,]+)", model):
            final_model = self.get_model(match.group(1))
            models = [(mn, self.get_model(mn)) for mn in match.group(2).split(",")]
            return StackingClassifier(models, final_estimator=final_model, n_jobs=self.n_jobs)
        raise ValueError()

    def __getitem__(self, model):
        return self.get_model(model)