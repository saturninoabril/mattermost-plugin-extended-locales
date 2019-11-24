import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {RECEIVED_EXTENDED_LANGUAGES} from './action_types';
import {id as pluginId} from './manifest';

export const getLocales = () => async (dispatch, getState) => {
    console.log('getLocales');
    fetch(getPluginServerRoute(getState()) + '/get_languages').then((r) => r.json()).then((r) => {
        const c = r.reduce((acc, p) => {
            acc[p.value] = p;
            return acc;
        }, {})
        console.log('plugin getLocales:', c);
        dispatch({
            type: RECEIVED_EXTENDED_LANGUAGES,
            data: c,
        });
    });
};

export const getPluginServerRoute = (state) => {
    const config = getConfig(state);

    let basePath = '/';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + pluginId;
};
