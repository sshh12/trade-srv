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
from sklearn.model_selection import cross_validate
from sklearn.pipeline import make_pipeline
from sklearn.preprocessing import StandardScaler
from sklearn.svm import SVC
from xgboost import XGBRegressor, XGBClassifier
from lightgbm import LGBMClassifier
from catboost import CatBoostClassifier
import numpy as np
import re


class ModelRepo:
    def __init__(self, seed=None, n_jobs=None):
        self.seed = seed
        self.n_jobs = n_jobs

    def get_model(self, name):
        raise NotImplementedError()

    def __getitem__(self, model):
        return self.get_model(model)

    def iter_fit(self, X, y=None, models=[]):
        for m in models:
            try:
                if isinstance(m, str):
                    model = self.get_model(m)
                    name = m
                else:
                    model = m
                    name = str(m)
                model.fit(X, y)
                yield name, model
            except Exception as e:
                print(e)
                yield name, None

    def iter_fit_cv(self, X, y, models=[], cv=3, scoring="accuracy", max_score=True):
        for m in models:
            try:
                if isinstance(m, str):
                    model = self.get_model(m)
                    name = m
                else:
                    model = m
                    name = str(m)
                result = cross_validate(model, X, y, cv=cv, return_estimator=True, scoring=scoring)
                if max_score:
                    clf = result["estimator"][np.argmax(result["test_score"])]
                else:
                    clf = result["estimator"][np.argmin(result["test_score"])]
                del result["estimator"]
                yield name, clf, result
            except Exception as e:
                print(e)
                yield name, None


class ClassifierRepo(ModelRepo):
    def get_model(self, model):
        if model == "ridgecv":
            return RidgeCV()
        elif model == "scaledsvc":
            return make_pipeline(StandardScaler(), SVC(gamma="auto"))
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
        elif match := re.match(r"gradboost(\d+)", model):
            return GradientBoostingClassifier(n_estimators=int(match.group(1)), random_state=self.seed)
        elif match := re.match(r"xgboost(\d+)", model):
            return XGBClassifier(n_estimators=int(match.group(1)), random_state=self.seed, use_label_encoder=False)
        elif match := re.match(r"lgbm(\d+)", model):
            return LGBMClassifier(
                n_estimators=int(match.group(1)),
                n_jobs=self.n_jobs,
                verbose=0,
                random_state=self.seed,
            )
        elif match := re.match(r"catboost(\d+)", model):
            return CatBoostClassifier(
                n_estimators=int(match.group(1)),
                thread_count=self.n_jobs,
                verbose=0,
                random_state=self.seed,
            )
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


class RegressionRepo(ModelRepo):
    def get_model(self, model):
        if match := re.match(r"ridgecv", model):
            return RidgeCV()
        elif match := re.match(r"rf(\d+)", model):
            return RandomForestRegressor(n_jobs=self.n_jobs, n_estimators=int(match.group(1)), random_state=self.seed)
        elif match := re.match(r"rf(\d+)d(\d+)", model):
            return RandomForestRegressor(
                n_jobs=self.n_jobs,
                n_estimators=int(match.group(1)),
                random_state=self.seed,
                max_depth=int(match.group(2)),
            )
        elif match := re.match(r"ada(\d+)", model):
            return AdaBoostRegressor(n_estimators=int(match.group(1)), random_state=self.seed)
        elif match := re.match(r"xtrees(\d+)", model):
            return ExtraTreesRegressor(n_estimators=int(match.group(1)), random_state=self.seed, n_jobs=self.n_jobs)
        elif match := re.match(r"gradboost(\d+)", model):
            return GradientBoostingRegressor(n_estimators=int(match.group(1)), random_state=self.seed)
        elif match := re.match(r"vote:([\w,]+)", model):
            models = [(mn, self.get_model(mn)) for mn in match.group(1).split(",")]
            return VotingRegressor(models, n_jobs=self.n_jobs)
        elif match := re.match(r"stack:(\w+):([\w,]+)", model):
            final_model = self.get_model(match.group(1))
            models = [(mn, self.get_model(mn)) for mn in match.group(2).split(",")]
            return StackingRegressor(models, final_estimator=final_model, n_jobs=self.n_jobs)
        raise ValueError()